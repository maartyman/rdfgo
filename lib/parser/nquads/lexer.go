package nquads

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Options holds the options for parsing.
type Options struct{}

// Parse parses the input stream and returns a channel of quads and an error channel.
func Parse(stream io.Reader, options Options) (chan interfaces.IQuad, chan error) {
	yyErrorVerbose = true
	errChan := make(chan error, 1)
	tokens := make(chan string, 100)
	out := make(chan interfaces.IQuad, 1)

	lex := &lexer{
		tokens:       tokens,
		output:       out,
		errChan:      errChan,
		channelsOpen: true,
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
	tokens       chan string
	output       chan interfaces.IQuad
	errChan      chan error
	lastTok      string
	channelsOpen bool
	mux          sync.Mutex
}

func unescapeLiteral(s string) (string, error) {
	// Wrap in double quotes so strconv.Unquote can decode it
	unquoted, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		return "", err
	}
	return unquoted, nil
}

func (l *lexer) Lex(lval *yySymType) int {
	tok, ok := <-l.tokens
	if !ok {
		// Channel is closed, signaling EOF
		return 0
	}
	l.lastTok = tok

	switch {
	case strings.HasPrefix(tok, "<") && strings.HasSuffix(tok, ">"):
		lval.term = NewNamedNode(tok)
		return _NNODE
	case strings.HasPrefix(tok, "_:"):
		lval.term = NewBlankNode(tok)
		return _BNODE
	case strings.HasPrefix(tok, "."):
		return _DOT
	case strings.HasPrefix(tok, "@"):
		lval.str = tok[1:]
		return _LANGTAG
	case strings.HasPrefix(tok, "^^"):
		lval.term = NewNamedNode(tok[2:])
		return _DATATYPE
	case strings.HasPrefix(tok, "\"") && strings.HasSuffix(tok, "\""):
		raw := tok[1 : len(tok)-1]
		unescaped, err := unescapeLiteral(raw)
		if err != nil {
			return _ERROR
		}
		lval.str = unescaped
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
