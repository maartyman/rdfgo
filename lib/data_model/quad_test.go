package rdfgo

import (
	"errors"
	"github.com/maartyman/rdfgo/interfaces"
	"testing"
)

func utilTermIsWrong(
	t *testing.T,
	returnedError error,
	q interfaces.IQuad,
	expectedError error,
	testErrorMessage string,
) {
	if returnedError == nil || !errors.Is(returnedError, expectedError) {
		if errors.Is(returnedError, expectedError) {
			t.Errorf("%s\n%s\n%s\ntrue", testErrorMessage, returnedError.Error(), expectedError.Error())
		} else {
			t.Errorf("%s\n%s\n%s\nfalse", testErrorMessage, returnedError.Error(), expectedError.Error())
		}
	}
	if q != nil {
		t.Errorf("Quad should be nil if error")
	}
}

func utilTermIsCorrect(t *testing.T, err error, q interfaces.IQuad, testErrorMessage string) {
	if err != nil {
		t.Errorf(testErrorMessage)
	}
	if q == nil {
		t.Errorf("Quad should not be nil if no error")
	}
}

func TestNewQuad_SubjectLimitations(t *testing.T) {
	l1, err1 := NewQuad(nil, NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l2, err2 := NewQuad(NewBlankNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l3, err3 := NewQuad(NewDefaultGraph(), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l4, err4 := NewQuad(NewLiteral("s", "", nil), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l5, err5 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l6, err6 := NewQuad(l5, NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l7, err7 := NewQuad(NewVariable("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))

	utilTermIsWrong(t, err1, l1, SubjectTermTypeError, "Quad subject cannot be nil")
	utilTermIsCorrect(t, err2, l2, "Quad subject can be a BlankNode")
	utilTermIsWrong(t, err3, l3, SubjectTermTypeError, "Quad subject cannot be a DefaultGraph")
	utilTermIsWrong(t, err4, l4, SubjectTermTypeError, "Quad subject cannot be a Literal")
	utilTermIsCorrect(t, err5, l5, "Quad subject can be a NamedNode")
	utilTermIsCorrect(t, err6, l6, "Quad subject can be a Quad")
	utilTermIsCorrect(t, err7, l7, "Quad subject can be a Variable")
}

func TestNewQuad_PredicateLimitations(t *testing.T) {
	l1, err1 := NewQuad(NewNamedNode("s"), nil, NewNamedNode("o"), NewNamedNode("g"))
	l2, err2 := NewQuad(NewNamedNode("s"), NewBlankNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l3, err3 := NewQuad(NewNamedNode("s"), NewDefaultGraph(), NewNamedNode("o"), NewNamedNode("g"))
	l4, err4 := NewQuad(NewNamedNode("s"), NewLiteral("p", "", nil), NewNamedNode("o"), NewNamedNode("g"))
	l5, err5 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l6, err6 := NewQuad(NewNamedNode("s"), l5, NewNamedNode("o"), NewNamedNode("g"))
	l7, err7 := NewQuad(NewNamedNode("s"), NewVariable("p"), NewNamedNode("o"), NewNamedNode("g"))

	utilTermIsWrong(t, err1, l1, PredicateTermTypeError, "Quad predicate cannot be nil")
	utilTermIsWrong(t, err2, l2, PredicateTermTypeError, "Quad predicate cannot be a BlankNode")
	utilTermIsWrong(t, err3, l3, PredicateTermTypeError, "Quad predicate cannot be a DefaultGraph")
	utilTermIsWrong(t, err4, l4, PredicateTermTypeError, "Quad predicate cannot be a Literal")
	utilTermIsCorrect(t, err5, l5, "Quad predicate can be a NamedNode")
	utilTermIsWrong(t, err6, l6, PredicateTermTypeError, "Quad predicate cannot be a Quad")
	utilTermIsCorrect(t, err7, l7, "Quad predicate can be a Variable")
}

func TestNewQuad_ObjectLimitations(t *testing.T) {
	l1, err1 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), nil, NewNamedNode("g"))
	l2, err2 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewBlankNode("o"), NewNamedNode("g"))
	l3, err3 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewDefaultGraph(), NewNamedNode("g"))
	l4, err4 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewLiteral("o", "", nil), NewNamedNode("g"))
	l5, err5 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l6, err6 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), l5, NewNamedNode("g"))
	l7, err7 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewVariable("o"), NewNamedNode("g"))

	utilTermIsWrong(t, err1, l1, ObjectTermTypeError, "Quad object cannot be nil")
	utilTermIsCorrect(t, err2, l2, "Quad object can be a BlankNode")
	utilTermIsWrong(t, err3, l3, ObjectTermTypeError, "Quad object cannot be a DefaultGraph")
	utilTermIsCorrect(t, err4, l4, "Quad object can be a Literal")
	utilTermIsCorrect(t, err5, l5, "Quad object can be a NamedNode")
	utilTermIsCorrect(t, err6, l6, "Quad object can be a Quad")
	utilTermIsCorrect(t, err7, l7, "Quad object can be a Variable")
}

func TestNewQuad_GraphLimitations(t *testing.T) {
	l1, err1 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), nil)
	l2, err2 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewBlankNode("g"))
	l3, err3 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewDefaultGraph())
	l4, err4 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewLiteral("g", "", nil))
	l5, err5 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l6, err6 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), l5)
	l7, err7 := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewVariable("g"))

	utilTermIsCorrect(t, err1, l1, "Quad graph can be nil")
	if !l1.GetGraph().Equals(NewDefaultGraph()) {
		t.Errorf("Quad graph should be a DefaultGraph if nil")
	}
	utilTermIsCorrect(t, err2, l2, "Quad graph can be a BlankNode")
	utilTermIsCorrect(t, err3, l3, "Quad graph can be a DefaultGraph")
	utilTermIsWrong(t, err4, l4, GraphTermTypeError, "Quad graph cannot be a Literal")
	utilTermIsCorrect(t, err5, l5, "Quad graph can be a NamedNode")
	utilTermIsWrong(t, err6, l6, GraphTermTypeError, "Quad graph cannot be a Quad")
	utilTermIsCorrect(t, err7, l7, "Quad graph can be a Variable")
}

func TestQuad_GetSubject(t *testing.T) {
	l1, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	if !l1.GetSubject().Equals(NewNamedNode("s")) {
		t.Errorf("Quad subject should equal s")
	}
}

func TestQuad_GetPredicate(t *testing.T) {
	l1, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	if !l1.GetPredicate().Equals(NewNamedNode("p")) {
		t.Errorf("Quad predicate should equal p")
	}
}

func TestQuad_GetObject(t *testing.T) {
	l1, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	if !l1.GetObject().Equals(NewNamedNode("o")) {
		t.Errorf("Quad object should equal o")
	}
}

func TestQuad_GetGraph(t *testing.T) {
	l1, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	if !l1.GetGraph().Equals(NewNamedNode("g")) {
		t.Errorf("Quad graph should equal g")
	}
}

func TestQuad_Equals(t *testing.T) {
	l1, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l2, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l3, _ := NewQuad(NewNamedNode("s1"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	l4, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p1"), NewNamedNode("o"), NewNamedNode("g"))
	l5, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o1"), NewNamedNode("g"))
	l6, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g1"))
	l7, _ := NewQuad(NewNamedNode("s1"), NewNamedNode("p1"), NewNamedNode("o1"), NewNamedNode("g1"))
	l8 := NewNamedNode("")
	if !l1.Equals(l1) {
		t.Errorf("A quad should equal itself")
	}
	if !l1.Equals(l2) {
		t.Errorf("A quad should equal a quad with the same values for all fields")
	}
	if l1.Equals(l3) {
		t.Errorf("A quad should not equal a quad with a different subject")
	}
	if l1.Equals(l4) {
		t.Errorf("A quad should not equal a quad with a different predicate")
	}
	if l1.Equals(l5) {
		t.Errorf("A quad should not equal a quad with a different object")
	}
	if l1.Equals(l6) {
		t.Errorf("A quad should not equal a quad with a different graph")
	}
	if l1.Equals(l7) {
		t.Errorf("A quad should not equal a quad with different values for all fields")
	}
	if l1.Equals(l8) {
		t.Errorf("A quad should not equal a NamedNode")
	}
}

func TestQuad_EqualsNil(t *testing.T) {
	l1, _ := NewQuad(NewNamedNode("s"), NewNamedNode("p"), NewNamedNode("o"), NewNamedNode("g"))
	if l1.Equals(nil) {
		t.Errorf("Quad should not equal nil")
	}
}

func TestQuad_ToString(t *testing.T) {
	tests := []struct {
		name     string
		quad     Quad
		expected string
	}{
		{
			name: "All components present",
			quad: Quad{
				subject:   NewNamedNode("<http://example.com/subject>"),
				predicate: NewNamedNode("<http://example.com/predicate>"),
				object:    NewNamedNode("<http://example.com/object>"),
				graph:     NewNamedNode("<http://example.com/graph>"),
			},
			expected: "<http://example.com/subject> <http://example.com/predicate> <http://example.com/object> <http://example.com/graph>",
		},
		{
			name: "Empty subject",
			quad: Quad{
				subject:   NewNamedNode(""),
				predicate: NewNamedNode("<http://example.com/predicate>"),
				object:    NewNamedNode("<http://example.com/object>"),
				graph:     NewNamedNode("<http://example.com/graph>"),
			},
			expected: "<> <http://example.com/predicate> <http://example.com/object> <http://example.com/graph>",
		},
		{
			name: "Empty predicate",
			quad: Quad{
				subject:   NewNamedNode("<http://example.com/subject>"),
				predicate: NewNamedNode(""),
				object:    NewNamedNode("<http://example.com/object>"),
				graph:     NewNamedNode("<http://example.com/graph>"),
			},
			expected: "<http://example.com/subject> <> <http://example.com/object> <http://example.com/graph>",
		},
		{
			name: "Empty object",
			quad: Quad{
				subject:   NewNamedNode("<http://example.com/subject>"),
				predicate: NewNamedNode("<http://example.com/predicate>"),
				object:    NewNamedNode(""),
				graph:     NewNamedNode("<http://example.com/graph>"),
			},
			expected: "<http://example.com/subject> <http://example.com/predicate> <> <http://example.com/graph>",
		},
		{
			name: "Empty graph",
			quad: Quad{
				subject:   NewNamedNode("<http://example.com/subject>"),
				predicate: NewNamedNode("<http://example.com/predicate>"),
				object:    NewNamedNode("<http://example.com/object>"),
				graph:     NewNamedNode(""),
			},
			expected: "<http://example.com/subject> <http://example.com/predicate> <http://example.com/object> <>",
		},
		{
			name: "All components empty",
			quad: Quad{
				subject:   NewNamedNode(""),
				predicate: NewNamedNode(""),
				object:    NewNamedNode(""),
				graph:     NewNamedNode(""),
			},
			expected: "<> <> <> <>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.quad.ToString()
			if actual != tt.expected {
				t.Errorf("ToString() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
