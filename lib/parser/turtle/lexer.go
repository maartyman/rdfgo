package turtle

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"io"
	"regexp"
	"strings"
	"sync"
)

// Options holds the options for parsing.
type Options struct {
	// BaseIRI is the base URI to use for resolving relative URIs.
	BaseIRI string
	// DataFactory is the data factory to use for creating RDF terms.
	DataFactory interfaces.IDataFactory
}

// Parse parses the input stream and returns a channel of quads and an error channel.
func Parse(stream io.Reader, options Options) (interfaces.IStream, chan error) {
	if options.DataFactory == nil {
		options.DataFactory = NewDataFactory()
	}
	yyErrorVerbose = true
	errChan := make(chan error, 1)
	tokens := make(chan token, 100)
	out := make(interfaces.IStream, 1)
	prefixes := make(map[string]string)

	lex := &lexer{
		tokens:       tokens,
		output:       out,
		errChan:      errChan,
		prefixes:     prefixes,
		base:         options.BaseIRI,
		channelsOpen: true,
		dataFactory:  options.DataFactory,
	}

	go func() {
		scanner := bufio.NewScanner(stream)
		defer close(tokens)
		for scanner.Scan() {
			lex.tokenizeLine(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			lex.mux.Lock()
			if lex.channelsOpen {
				errChan <- err
				lex.channelsOpen = false
				close(out)
				close(errChan)
			}
			lex.mux.Unlock()
		}
	}()

	go func() {
		_ = yyParse(lex)
		lex.mux.Lock()
		if lex.channelsOpen {
			lex.channelsOpen = false
			close(out)
			close(errChan)
		}
		lex.mux.Unlock()
	}()
	return out, errChan
}

type lexer struct {
	tokens           chan token
	output           interfaces.IStream
	errChan          chan error
	prefixes         map[string]string
	lastTok          string
	base             string
	multiLineLiteral string
	channelsOpen     bool
	mux              sync.Mutex
	dataFactory      interfaces.IDataFactory
}

type token struct {
	Type int
	Val  string
}

func (l *lexer) Lex(lval *yySymType) int {
	tok, ok := <-l.tokens
	if !ok {
		return 0
	}
	l.lastTok = tok.Val

	switch tok.Type {
	case _NUMBER:
		var datatypeIRI interfaces.INamedNode
		if strings.ContainsAny(tok.Val, "eE") {
			datatypeIRI = IRI.XSD.Double
		} else if strings.Contains(tok.Val, ".") {
			datatypeIRI = IRI.XSD.Decimal
		} else {
			datatypeIRI = IRI.XSD.Integer
		}
		lval.literal = l.dataFactory.Literal(tok.Val, "", datatypeIRI)
	case _LITERAL:
		// parse literal, lang and datatype
		lval.str = unescape(tok.Val[1 : len(tok.Val)-1])
	case _PNAME, _PVALUE, _NNODE, _BNODE, _LANGTAG:
		lval.str = unescape(tok.Val)
	}
	return tok.Type
}

func (l *lexer) Error(e string) {
	msg := fmt.Sprintf("Syntax error near token: %q (%s)", l.lastTok, e)

	l.mux.Lock()
	if l.channelsOpen {
		l.errChan <- errors.New(msg)
		l.channelsOpen = false
		close(l.output)
		close(l.errChan)
	}
	l.mux.Unlock()
}

func unescape(s string) string {
	s = replacer.Replace(s)

	s = reUnicode.ReplaceAllStringFunc(s, func(m string) string {
		var code int
		// We ignore the errors as the regex will only match valid Unicode escape sequences
		if strings.HasPrefix(m, `\u`) {
			_, _ = fmt.Sscanf(m, `\u%04x`, &code)
		} else {
			_, _ = fmt.Sscanf(m, `\U%08x`, &code)
		}
		return string(rune(code))
	})

	return s
}

func processPrefixed(str string, tokenChan chan token) {
	index := strings.Index(str, ":")
	if index == -1 {
		tokenChan <- token{_ERROR, str}
		return
	}
	tokenChan <- token{_PNAME, str[:index]}
	if index != len(str)-1 {
		if str[len(str)-1] == '.' {
			tokenChan <- token{_PVALUE, str[index+1 : len(str)-1]}
			tokenChan <- token{_DOT, "."}
		} else {
			tokenChan <- token{_PVALUE, str[index+1:]}
		}
	}
}

func (l *lexer) tokenizeLine(line string) {
	if l.multiLineLiteral != "" {
		line = l.multiLineLiteral + "\n" + line
	}
	rawTokens := reTokenize.FindAllString(line, -1)

	for _, rt := range rawTokens {
		firstChar := rt[0]
		switch {
		case rt == ".":
			l.tokens <- token{_DOT, rt}
		case rt == ";":
			l.tokens <- token{_SEMICOLON, rt}
		case rt == ",":
			l.tokens <- token{_COMMA, rt}
		case rt == "a":
			l.tokens <- token{_A, rt}
		case rt == "@prefix":
			l.tokens <- token{_PREFIX, rt}
		case strings.ToLower(rt) == "prefix":
			l.tokens <- token{_SPARQLPREFIX, rt}
		case rt == "@base":
			l.tokens <- token{_BASE, rt}
		case strings.ToLower(rt) == "base":
			l.tokens <- token{_SPARQLBASE, rt}
		case rt == "[":
			l.tokens <- token{_BRACKETOPEN, rt}
		case rt == "]":
			l.tokens <- token{_BRACKETCLOSE, rt}
		case rt == "(":
			l.tokens <- token{_LPAREN, rt}
		case rt == ")":
			l.tokens <- token{_RPAREN, rt}
		case rt == "true":
			l.tokens <- token{_TRUE, rt}
		case rt == "false":
			l.tokens <- token{_FALSE, rt}
		case strings.HasPrefix(rt, `^^`):
			l.tokens <- token{_DATATYPE, rt}
			if rt[2] == '<' {
				l.tokens <- token{_NNODE, rt[3 : len(rt)-1]}
			} else {
				processPrefixed(rt[2:], l.tokens)
			}
		case strings.HasPrefix(rt, `@`):
			l.tokens <- token{_LANGTAG, rt[1:]}
		case strings.HasPrefix(rt, "<") && strings.HasSuffix(rt, ">"):
			l.tokens <- token{_NNODE, rt[1 : len(rt)-1]}
		case strings.HasPrefix(rt, "_:"):
			l.tokens <- token{_BNODE, rt[2:]}
		case strings.HasPrefix(rt, `"""`) && strings.HasSuffix(rt, `"""`) && len(rt) > 4:
			l.multiLineLiteral = ""
			l.tokens <- token{_LITERAL, rt[2 : len(rt)-2]}
		case strings.HasPrefix(rt, `'''`) && strings.HasSuffix(rt, `'''`) && len(rt) > 4:
			l.multiLineLiteral = ""
			l.tokens <- token{_LITERAL, rt[2 : len(rt)-2]}
		case strings.HasPrefix(rt, `"""`):
			l.multiLineLiteral = rt
		case strings.HasPrefix(rt, `'''`):
			l.multiLineLiteral = rt
		case strings.HasPrefix(rt, "\""):
			l.tokens <- token{_LITERAL, rt}
		case strings.HasPrefix(rt, "'"):
			l.tokens <- token{_LITERAL, rt}
		case strings.HasPrefix(rt, "#"):
			// ignore comments
			continue
		case (firstChar >= '0' && firstChar <= '9') || firstChar == '-' || firstChar == '+' || firstChar == '.':
			l.tokens <- token{_NUMBER, rt}
		default:
			processPrefixed(rt, l.tokens)
		}
	}
}

// Regular expressions
var (
	numericPattern                     = `[+-]?(?:(?:[0-9]+(?:\.[0-9]+)?)|(?:\.[0-9]+))(?:\.?[eE][+-]?[0-9]+)?`
	booleanPattern                     = `\btrue\b|\bfalse\b`
	iriPattern                         = `<[^>]*>`
	fullMultilineDoubleQuoteLiteral    = `"""(?:(?s).*?)[^\\]?"""`
	fullMultilineSingleQuoteLiteral    = `'''(?:(?s).*?)[^\\]?'''`
	partialMultilineDoubleQuoteLiteral = `"""(?:(?s).*)`
	partialMultilineSingleQuoteLiteral = `'''(?:(?s).*)`
	literalPatternDouble               = `"(?:[^"\\]|\\.)*"`
	literalPatternSingle               = `'(?:[^'\\]|\\.)*'`
	datatypePattern                    = `\^\^` + iriPattern
	langTagPattern                     = `@[a-zA-Z]+(?:-[a-zA-Z0-9]+)*`
	blankNodePattern                   = `_:[\S]+`
	symbolPattern                      = `[\.,;\[\]\(\)]`
	prefixedNamePattern                = `(?:[^\(\)\s,;\[\]@\\#])*:(?:\\[_~\.\-!$&'\(\)*+,;=/?#@%]|[^\(\)\s,;\[\]@\\#])*`
	keywordsPattern                    = `(?i)prefix|(?i)base|@prefix|@base|a`
	commentPattern                     = `#.*`
	fallbackPattern                    = `[^\s]+`

	combinedPattern = strings.Join([]string{
		numericPattern,
		booleanPattern,
		iriPattern,
		fullMultilineDoubleQuoteLiteral,
		fullMultilineSingleQuoteLiteral,
		partialMultilineDoubleQuoteLiteral,
		partialMultilineSingleQuoteLiteral,
		literalPatternDouble,
		literalPatternSingle,
		datatypePattern,
		langTagPattern,
		blankNodePattern,
		symbolPattern,
		prefixedNamePattern,
		keywordsPattern,
		commentPattern,
		fallbackPattern,
	}, "|")

	reTokenize = regexp.MustCompile(combinedPattern)
	reUnicode  = regexp.MustCompile(`\\u([0-9A-Fa-f]{4})|\\U([0-9A-Fa-f]{8})`)
	replacer   = strings.NewReplacer(
		`\\`, `\`,
		`\"`, `"`,
		`\'`, `'`,
		`\n`, "\n",
		`\t`, "\t",
		`\r`, "\r",
		`\b`, "\b",
		`\f`, "\f",
		`\_`, "_",
		`\~`, "~",
		`\.`, ".",
		`\-`, "-",
		`\!`, "!",
		`\$`, "$",
		`\&`, "&",
		`\(`, "(",
		`\)`, ")",
		`\*`, "*",
		`\+`, "+",
		`\,`, ",",
		`\;`, ";",
		`\=`, "=",
		`\/`, "/",
		`\?`, "?",
		`\#`, "#",
		`\@`, "@",
		`\%`, "%",
	)
)
