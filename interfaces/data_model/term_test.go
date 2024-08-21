package interfaces

import "testing"

func TestTermType_String(t *testing.T) {
	if NamedNodeType.String() != namedNodeTypeString {
		t.Errorf("NamedNodeType string should be %s", namedNodeTypeString)
	}
	if LiteralType.String() != literalTypeString {
		t.Errorf("LiteralType string should be %s", literalTypeString)
	}
	if BlankNodeType.String() != blankNodeTypeString {
		t.Errorf("BlankNodeType string should be %s", blankNodeTypeString)
	}
	if VariableType.String() != variableTypeString {
		t.Errorf("VariableType string should be %s", variableTypeString)
	}
	if DefaultGraphType.String() != defaultGraphTypeString {
		t.Errorf("DefaultGraphType string should be %s", defaultGraphTypeString)
	}
	if QuadType.String() != quadTypeString {
		t.Errorf("QuadType string should be %s", quadTypeString)
	}
}

func TestTermType_EnumIndex(t *testing.T) {
	if NamedNodeType.EnumIndex() != 0 {
		t.Errorf("NamedNodeType enum index should be 0")
	}
	if LiteralType.EnumIndex() != 1 {
		t.Errorf("LiteralType enum index should be 1")
	}
	if BlankNodeType.EnumIndex() != 2 {
		t.Errorf("BlankNodeType enum index should be 2")
	}
	if VariableType.EnumIndex() != 3 {
		t.Errorf("VariableType enum index should be 3")
	}
	if DefaultGraphType.EnumIndex() != 4 {
		t.Errorf("DefaultGraphType enum index should be 4")
	}
	if QuadType.EnumIndex() != 5 {
		t.Errorf("QuadType enum index should be 5")
	}
}
