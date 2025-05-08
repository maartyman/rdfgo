package nquads

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
	tokens := make(chan string, 100)
	out := make(interfaces.IStream, 1)

	lex := &lexer{
		tokens:         tokens,
		output:         out,
		errChan:        errChan,
		channelsOpen:   true,
		dataFactory:    options.DataFactory,
		blankNodeIndex: make(map[string]interfaces.IBlankNode),
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
	tokens         chan string
	output         interfaces.IStream
	errChan        chan error
	lastTok        string
	channelsOpen   bool
	mux            sync.Mutex
	dataFactory    interfaces.IDataFactory
	blankNodeIndex map[string]interfaces.IBlankNode
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

var reUnicode = regexp.MustCompile(`\\u([0-9A-Fa-f]{4})|\\U([0-9A-Fa-f]{8})`)
var replacer = strings.NewReplacer(
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

func (l *lexer) Lex(lval *yySymType) int {
	tok, ok := <-l.tokens
	if !ok {
		// Channel is closed, signaling EOF
		return 0
	}
	l.lastTok = tok

	switch {
	case strings.HasPrefix(tok, "<") && strings.HasSuffix(tok, ">"):
		lval.term = l.dataFactory.NamedNode(unescape(tok[1 : len(tok)-1]))
		return _NNODE
	case strings.HasPrefix(tok, "_:"):
		// Check if the blank node is already created
		if node, ok := l.blankNodeIndex[tok]; ok {
			lval.term = node
			return _BNODE
		}
		// Create a new blank node
		bn := l.dataFactory.BlankNode("")
		l.blankNodeIndex[tok] = bn
		lval.term = bn
		return _BNODE
	case strings.HasPrefix(tok, "."):
		return _DOT
	case strings.HasPrefix(tok, "@"):
		lval.str = tok[1:]
		return _LANGTAG
	case strings.HasPrefix(tok, "^^"):
		lval.term = l.dataFactory.NamedNode(tok[2:])
		return _DATATYPE
	case strings.HasPrefix(tok, "\"") && strings.HasSuffix(tok, "\""):
		raw := tok[1 : len(tok)-1]
		lval.str = unescape(raw)
		return _LITERALVALUE
	default:
		return _ERROR
	}
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

var reIri = `<(?:[^>]+)>`
var reLiteral = `"(?:(?:[^"\\]|\\.)*)"`
var reLang = `(?:@[a-zA-Z]+(?:-[a-zA-Z0-9]+)*)`
var reDatatype = `\^\^` + reIri
var reBlankNode = `_:(?:[A-Za-z0-9_:.]*[A-Za-z0-9_:])`
var reDot = `\.`
var reFallback = `\S+`

var re = regexp.MustCompile(
	reIri + `|` + reLiteral + `|` + reLang + `|` + reDatatype + `|` + reBlankNode + `|` + reDot + `|` + reFallback,
)

func stripComments(line string) string {
	inIRI := false
	inLiteral := false
	escaped := false

	for i := 0; i < len(line); i++ {
		c := line[i]

		if inLiteral {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
			} else if c == '"' {
				inLiteral = false
			}
			continue
		}

		if inIRI {
			if c == '>' {
				inIRI = false
			}
			continue
		}

		switch c {
		case '<':
			inIRI = true
		case '"':
			inLiteral = true
		case '#':
			// Not inside IRI or literal => it's a comment
			return line[:i]
		}
	}

	return line
}

func (l *lexer) tokenizeLine(line string) {
	line = stripComments(line)
	rawTokens := re.FindAllString(line, -1)

	for _, rt := range rawTokens {
		l.tokens <- rt
	}
}
