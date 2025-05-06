package nquads_parser

import (
	"bufio"
	"errors"
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func ParseNQuads(stream io.Reader) (chan interfaces.IQuad, chan error) {
	errChan := make(chan error, 1)
	tokens := make(chan string)
	out := make(chan interfaces.IQuad)

	go func() {
		scanner := bufio.NewScanner(stream)
		defer close(tokens)
		for scanner.Scan() {
			text := scanner.Text()
			tokensInLine := tokenizeLine(text)
			for _, tok := range tokensInLine {
				tokens <- tok
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
		}
	}()

	lexer := &Lexer{
		tokens:  tokens,
		output:  out,
		errChan: errChan,
	}

	go func() {
		code := yyParse(lexer)
		if code != 0 {
			// an error occurred during parsing, the output channel and error channel are already closed
			return
		}
		close(out)
		close(errChan)
	}()
	return out, errChan
}

type Lexer struct {
	tokens  chan string
	output  chan interfaces.IQuad
	errChan chan error
}

func unescapeLiteral(s string) (string, error) {
	// Wrap in double quotes so strconv.Unquote can decode it
	unquoted, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		return "", err
	}
	return unquoted, nil
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok, ok := <-l.tokens
	if !ok {
		// Channel is closed, signaling EOF
		return 0
	}

	switch {
	case strings.HasPrefix(tok, "<") && strings.HasSuffix(tok, ">"):
		lval.term = NewNamedNode(tok)
		return NNODE
	case strings.HasPrefix(tok, "_:"):
		lval.term = NewBlankNode(tok)
		return BNODE
	case strings.HasPrefix(tok, "."):
		return DOT
	case strings.HasPrefix(tok, "@"):
		lval.str = tok[1:]
		return LANGTAG
	case strings.HasPrefix(tok, "^^"):
		lval.term = NewNamedNode(tok[2:])
		return DATATYPE
	case strings.HasPrefix(tok, "\"") && strings.HasSuffix(tok, "\""):
		raw := tok[1 : len(tok)-1]
		unescaped, err := unescapeLiteral(raw)
		if err != nil {
			return ERROR
		}
		lval.str = unescaped
		return LITERALVALUE
	default:
		return ERROR
	}
}

func (l *Lexer) Error(e string) {
	l.errChan <- errors.New(e)
	close(l.output)
	close(l.errChan)
}

var reIri = `<([^>]+)>`
var reLiteral = `"((?:[^"\\]|\\.)*)"`
var reLang = `(@[a-zA-Z]+(?:-[a-zA-Z0-9]+)*)`
var reDatatype = `(\^\^` + reIri + `)`
var reBlankNode = `_:([A-Za-z0-9_:.]*[A-Za-z0-9_:])`
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

func tokenizeLine(line string) []string {
	line = stripComments(line)
	return re.FindAllString(line, -1)
}
