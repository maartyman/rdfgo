package rdfgo

import "testing"

func TestVariable_NewVariable(t *testing.T) {
	v1 := NewVariable("?v1")
	v2 := NewVariable("v1")
	if v1.GetValue() != "v1" {
		t.Errorf("Variable to string should equal ?v1")
	}
	if !v1.Equals(v2) {
		t.Errorf("Variable with and without `?` should be equal")
	}
}

func TestVariable_EqualsNil(t *testing.T) {
	v1 := NewVariable("v1")
	if v1.Equals(nil) {
		t.Errorf("Variable should not equal nil")
	}
}

func TestVariable_ToString(t *testing.T) {
	v1 := NewVariable("v1")
	if v1.ToString() != "?v1" {
		t.Errorf("Variable to string should equal ?v1")
	}
}
