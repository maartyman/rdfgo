package turtle_parser

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

type TurtleParserOptions struct {
	// Base is the base URI to use for resolving relative URIs.
	Base string
}

func ParseTurtle(stream io.Reader, opt *TurtleParserOptions) (chan interfaces.IQuad, chan error) {
	if opt == nil {
		opt = &TurtleParserOptions{}
	}
	yyErrorVerbose = true
	errChan := make(chan error, 1)
	tokens := make(chan Token, 1000)
	out := make(chan interfaces.IQuad, 1000)
	prefixes := make(map[string]string)

	lexer := &Lexer{
		tokens:       tokens,
		output:       out,
		errChan:      errChan,
		prefixes:     prefixes,
		base:         opt.Base,
		channelsOpen: true,
	}

	go func() {
		scanner := bufio.NewScanner(stream)
		defer close(tokens)
		for scanner.Scan() {
			lexer.tokenizeLine(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			lexer.mux.Lock()
			if lexer.channelsOpen {
				errChan <- err
				lexer.channelsOpen = false
				close(out)
				close(errChan)
			}
			lexer.mux.Unlock()
		}
	}()

	go func() {
		_ = yyParse(lexer)
		lexer.mux.Lock()
		if lexer.channelsOpen {
			lexer.channelsOpen = false
			close(out)
			close(errChan)
		}
		lexer.mux.Unlock()
	}()
	return out, errChan
}

type Lexer struct {
	tokens           chan Token
	output           chan interfaces.IQuad
	errChan          chan error
	prefixes         map[string]string
	lastTok          string
	base             string
	blankNodeCounter int
	multiLineLiteral string
	channelsOpen     bool
	mux              sync.Mutex
}

type Token struct {
	Type int
	Val  string
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok, ok := <-l.tokens
	if !ok {
		return 0
	}
	l.lastTok = tok.Val

	switch tok.Type {
	case NUMBER:
		var datatypeIRI string
		if strings.ContainsAny(tok.Val, "eE") {
			datatypeIRI = "http://www.w3.org/2001/XMLSchema#double"
		} else if strings.Contains(tok.Val, ".") {
			datatypeIRI = "http://www.w3.org/2001/XMLSchema#decimal"
		} else {
			datatypeIRI = "http://www.w3.org/2001/XMLSchema#integer"
		}
		lval.literal = NewLiteral(tok.Val, "", NewNamedNode(datatypeIRI))
	case LITERAL:
		// parse literal, lang and datatype
		unescaped := unescapeLiteral(tok.Val[1 : len(tok.Val)-1])
		lval.str = unescaped
	case PNAME, PVALUE, NNODE, BNODE, LANGTAG, DATATYPE:
		lval.str = tok.Val
	}
	return tok.Type
}

func (l *Lexer) Error(e string) {
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

func unescapeLiteral(s string) string {
	replacer := strings.NewReplacer(
		`\\`, `\`,
		`\"`, `"`,
		`\'`, `'`,
		`\n`, "\n",
		`\t`, "\t",
		`\r`, "\r",
		`\b`, "\b",
		`\f`, "\f",
	)
	s = replacer.Replace(s)

	// Handle \uXXXX and \UXXXXXXXX
	reUnicode := regexp.MustCompile(`\\u([0-9A-Fa-f]{4})|\\U([0-9A-Fa-f]{8})`)
	s = reUnicode.ReplaceAllStringFunc(s, func(m string) string {
		var code int
		if strings.HasPrefix(m, `\u`) {
			fmt.Sscanf(m, `\u%04x`, &code)
		} else {
			fmt.Sscanf(m, `\U%08x`, &code)
		}
		return string(rune(code))
	})

	return s
}

func (l *Lexer) newBlankNode() interfaces.IBlankNode {
	id := fmt.Sprintf("b%d", l.blankNodeCounter)
	l.blankNodeCounter++
	return NewBlankNode(id)
}

func (l *Lexer) tokenizeLine(line string) {
	if l.multiLineLiteral != "" {
		line = l.multiLineLiteral + "\n" + line
	}
	rawTokens := reTokenize.FindAllString(line, -1)

	for _, rt := range rawTokens {
		firstChar := rt[0]
		switch {
		case rt == ".":
			l.tokens <- Token{DOT, rt}
		case rt == ";":
			l.tokens <- Token{SEMICOLON, rt}
		case rt == ",":
			l.tokens <- Token{COMMA, rt}
		case rt == "a":
			l.tokens <- Token{A, rt}
		case rt == "@prefix":
			l.tokens <- Token{PREFIX, rt}
		case strings.ToLower(rt) == "prefix":
			l.tokens <- Token{SPARQLPREFIX, rt}
		case rt == "@base":
			l.tokens <- Token{BASE, rt}
		case strings.ToLower(rt) == "base":
			l.tokens <- Token{SPARQLBASE, rt}
		case rt == "[":
			l.tokens <- Token{BRACKETOPEN, rt}
		case rt == "]":
			l.tokens <- Token{BRACKETCLOSE, rt}
		case rt == "(":
			l.tokens <- Token{LPAREN, rt}
		case rt == ")":
			l.tokens <- Token{RPAREN, rt}
		case rt == "true":
			l.tokens <- Token{TRUE, rt}
		case rt == "false":
			l.tokens <- Token{FALSE, rt}
		case strings.HasPrefix(rt, `^^`):
			l.tokens <- Token{DATATYPE, rt[2:]}
		case strings.HasPrefix(rt, `@`):
			l.tokens <- Token{LANGTAG, rt[1:]}
		case strings.HasPrefix(rt, "<") && strings.HasSuffix(rt, ">"):
			l.tokens <- Token{NNODE, rt[1 : len(rt)-1]}
		case strings.HasPrefix(rt, "_:"):
			l.tokens <- Token{BNODE, rt[2:]}
		case strings.HasPrefix(rt, `"""`) && strings.HasSuffix(rt, `"""`) && len(rt) > 4:
			l.multiLineLiteral = ""
			l.tokens <- Token{LITERAL, rt[2 : len(rt)-2]}
		case strings.HasPrefix(rt, `'''`) && strings.HasSuffix(rt, `'''`) && len(rt) > 4:
			l.multiLineLiteral = ""
			l.tokens <- Token{LITERAL, rt[2 : len(rt)-2]}
		case strings.HasPrefix(rt, `"""`):
			l.multiLineLiteral = rt
		case strings.HasPrefix(rt, `'''`):
			l.multiLineLiteral = rt
		case strings.HasPrefix(rt, "\""):
			l.tokens <- Token{LITERAL, rt}
		case strings.HasPrefix(rt, "'"):
			l.tokens <- Token{LITERAL, rt}
		case strings.HasPrefix(rt, "#"):
			// ignore comments
			continue
		case (firstChar >= '0' && firstChar <= '9') || firstChar == '-' || firstChar == '+' || firstChar == '.':
			l.tokens <- Token{NUMBER, rt}
		case strings.Contains(rt, ":"):
			index := strings.Index(rt, ":")
			if index > 0 && rt[index-1] == '.' {
				l.tokens <- Token{ERROR, rt}
			}
			l.tokens <- Token{PNAME, rt[:index]}
			if index != len(rt)-1 {
				if rt[len(rt)-1] == '.' {
					l.tokens <- Token{PVALUE, rt[index+1 : len(rt)-1]}
					l.tokens <- Token{DOT, "."}
				} else {
					l.tokens <- Token{PVALUE, rt[index+1:]}
				}
			}
		default:
			// fallback if needed
			l.tokens <- Token{ERROR, rt}
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
	literalPatternDouble               = `"(?:(?:[^"\\]|\\.)*)"`
	literalPatternSingle               = `'(?:(?:[^'\\]|\\.)*)'`
	datatypePattern                    = `\^\^` + iriPattern
	langTagPattern                     = `@[a-zA-Z]+(?:-[a-zA-Z0-9]+)*`
	blankNodePattern                   = `_:[\S]+`
	symbolPattern                      = `[\.,;\[\]\(\)]`
	prefixedNamePattern                = `(?:[^\(\)\s,;\[\]@\\])*:(?:\\[_~\.\-!$&'\(\)*+,;=/?#@%]|[^\(\)\s,;\[\]@\\])*`
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
)
