package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	. "github.com/maartyman/rdfgo/lib/data_model"
	"strconv"
	"sync"
	"testing"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Error("Store should not be nil")
	}
}

func TestStoreImport(t *testing.T) {
	store := NewStore()
	quad, _ := NewQuad(
		NewNamedNode("http://example.com/s"),
		NewNamedNode("http://example.com/p"),
		NewNamedNode("http://example.com/o"),
		NewDefaultGraph(),
	)
	stream := ArrayToStream([]interfaces.IQuad{
		quad,
	})
	store.Import(interfaces.IStream(stream))
	if store.Size() <= 0 {
		t.Error("Store should have quads")
	}
}

func TestAddQuadFromTerms_NilTerms(t *testing.T) {
	store := NewStore()

	result := store.AddQuadFromTerms(nil, nil, nil, nil)
	if result != false {
		t.Errorf("Expected AddQuadFromTerms to return false when all terms are nil")
	}
}

func TestAddQuadFromTerms_FaultyQuad(t *testing.T) {
	store := NewStore()

	subject := NewLiteral("subject", "", NewNamedNode(""))
	predicate := NewNamedNode("predicate")
	object := NewNamedNode("object")
	graph := NewNamedNode("graph")

	_, err := NewQuad(subject, predicate, object, graph)
	if err == nil {
		t.Errorf("Expected NewQuad to return an error")
	}

	result := store.AddQuadFromTerms(subject, predicate, object, graph)
	if result != false {
		t.Errorf("Expected AddQuadFromTerms to return false when NewQuad would return an error")
	}
}

func TestAddQuadFromTerms_NilGraph(t *testing.T) {
	store := NewStore()

	result := store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		nil,
	)
	if result != true {
		t.Errorf("Expected AddQuadFromTerms to succeed even when the graph is nil")
	}
}

func TestAddQuadFromTerms_WithVariable(t *testing.T) {
	store := NewStore()

	result := store.AddQuadFromTerms(
		NewVariable("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	if result == true {
		t.Errorf("Expected AddQuadFromTerms to return false when the subject is a variable")
	}

	result = store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewVariable("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	if result == true {
		t.Errorf("Expected AddQuadFromTerms to return false when the subject is a variable")
	}

	result = store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewVariable("object"),
		NewNamedNode("graph"),
	)
	if result == true {
		t.Errorf("Expected AddQuadFromTerms to return false when the subject is a variable")
	}

	result = store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewVariable("graph"),
	)
	if result == true {
		t.Errorf("Expected AddQuadFromTerms to return false when the subject is a variable")
	}
}

func TestAddQuad_Duplicate(t *testing.T) {
	store := NewStore()

	store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	result := store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)

	if result != false {
		t.Errorf("Expected AddQuadFromTerms to return false when adding a duplicate quad")
	}

	if store.Size() != 1 {
		t.Errorf("Expected store size to be 1, but got %d", store.Size())
	}
}

func TestAddQuad_DuplicateWithDifferentGraph(t *testing.T) {
	store := NewStore()

	store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		nil,
	)
	result := store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)

	if result == false {
		t.Errorf(
			"Expected AddQuadFromTerms to return true when adding a duplicate quad with different graphs",
		)
	}

	if store.Size() != 2 {
		t.Errorf("Expected store size to be 1, but got %d", store.Size())
	}
}

func TestRemoveQuad_NotExist(t *testing.T) {
	store := NewStore()

	quad, _ := NewQuad(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	store.RemoveQuad(quad)

	if store.Size() != 0 {
		t.Errorf("Expected store size to remain 0 after trying to remove a non-existent quad")
	}
}

func TestRemoveQuad_Exist(t *testing.T) {
	store := NewStore()

	quad, _ := NewQuad(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	store.AddQuad(quad)
	store.RemoveQuad(quad)

	if store.Size() != 0 {
		t.Errorf("Expected store size to be 0 after removing a quad")
	}
}

func TestRemoveQuad_ExistDifferentGraphs(t *testing.T) {
	store := NewStore()

	quad1, _ := NewQuad(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph1"),
	)
	quad2, _ := NewQuad(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph2"),
	)
	store.AddQuad(quad1)
	store.AddQuad(quad2)
	store.RemoveQuad(quad1)

	if store.Has(quad1) {
		t.Errorf("Expected store to not have quad1 after removing it")
	}
	if store.Has(quad2) == false {
		t.Errorf("Expected store to have quad2 after removing quad1")
	}
	if store.Size() != 1 {
		t.Errorf("Expected store size to be 1 after removing a quad, but got %d", store.Size())
	}
}

func TestStore_Match(t *testing.T) {
	store := NewStore()

	store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewDefaultGraph(),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject1"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject2"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject3"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject4"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject1"),
		NewNamedNode("predicate1"),
		NewNamedNode("object1"),
		NewNamedNode("graph1"),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject1"),
		NewNamedNode("predicate1"),
		NewNamedNode("object2"),
		NewNamedNode("graph1"),
	)

	tests := []struct {
		name      string
		subject   interfaces.INamedNode
		predicate interfaces.INamedNode
		object    interfaces.INamedNode
		graph     interfaces.INamedNode
		expected  int
		message   string
		store     IStore
	}{
		{
			"Match_AllNil",
			nil,
			nil,
			nil,
			nil,
			8,
			"Expected Match to return all quads when all terms are nil",
			store,
		},
		{
			"Match_WrongTerms",
			NewLiteral("subject", "", NewNamedNode("")),
			nil,
			nil,
			nil,
			0,
			"Expected Match to return 0 quads when an impossible pattern has been given",
			store,
		},
		{
			"Match_AllWrongTerms",
			NewDefaultGraph(),
			NewDefaultGraph(),
			NewDefaultGraph(),
			NewLiteral("subject", "", NewNamedNode("")),
			0,
			"Expected Match to return 0 quads when an impossible pattern has been given",
			store,
		},

		{
			"Match_NoMatchSubject",
			NewNamedNode("no_match"),
			nil,
			nil,
			nil,
			0,
			"Expected Match to return 0 quads when there are no matches for subject",
			store,
		},
		{
			"Match_NoMatchPredicate",
			nil,
			NewNamedNode("no_match"),
			nil,
			nil,
			0,
			"Expected Match to return 0 quads when there are no matches for predicate",
			store,
		},
		{
			"Match_NoMatchObject",
			nil,
			nil,
			NewNamedNode("no_match"),
			nil,
			0,
			"Expected Match to return 0 quads when there are no matches for object",
			store,
		},
		{
			"Match_NoMatchGraph",
			nil,
			nil,
			nil,
			NewNamedNode("no_match"),
			0,
			"Expected Match to return 0 quads when there are no matches for graph",
			store,
		},

		{
			"Match_Subject",
			NewNamedNode("subject"),
			nil,
			nil,
			nil,
			2,
			"Expected Match to return 2 quads when matching subject",
			store,
		},
		{
			"Match_Predicate",
			nil,
			NewNamedNode("predicate"),
			nil,
			nil,
			6,
			"Expected Match to return 6 quads when matching predicate",
			store,
		},
		{
			"Match_Object",
			nil,
			nil,
			NewNamedNode("object"),
			nil,
			6,
			"Expected Match to return 6 quads when matching object",
			store,
		},
		{
			"Match_Graph",
			nil,
			nil,
			nil,
			NewNamedNode("graph"),
			5,
			"Expected Match to return 5 quads when matching graph",
			store,
		},
		{
			"Match_VariableSubject",
			NewVariable("s"),
			nil,
			nil,
			nil,
			8,
			"Expected Match to ignore Subject Variables and return 8 quads",
			store,
		},
		{
			"Match_VariablePredicate",
			nil,
			NewVariable("p"),
			nil,
			nil,
			8,
			"Expected Match to ignore Predicate Variables and return 8 quads",
			store,
		},
		{
			"Match_VariableObject",
			nil,
			nil,
			NewVariable("o"),
			nil,
			8,
			"Expected Match to ignore Object Variables and return 8 quads",
			store,
		},
		{
			"Match_VariableGraph",
			nil,
			nil,
			nil,
			NewVariable("g"),
			8,
			"Expected Match to ignore Graph Variables and return 8 quads",
			store,
		},

		{
			"Match_SubjectAndPredicate",
			NewNamedNode("subject1"),
			NewNamedNode("predicate1"),
			nil,
			nil,
			2,
			"Expected Match to return 2 quads when matching subject and predicate",
			store,
		},
		{
			"Match_SubjectPredicateObject",
			NewNamedNode("subject1"),
			NewNamedNode("predicate1"),
			NewNamedNode("object1"),
			nil,
			1,
			"Expected Match to return 1 quad when matching subject, predicate, and object",
			store,
		},
		{
			"Match_AllTerms",
			NewNamedNode("subject1"),
			NewNamedNode("predicate1"),
			NewNamedNode("object1"),
			NewNamedNode("graph1"),
			1,
			"Expected Match to return 1 quad when matching subject, predicate, object, and graph",
			store,
		},

		{
			"Match_EmptyStore",
			nil,
			nil,
			nil,
			nil,
			0,
			"Expected Match to return 0 quads when the store is empty",
			NewStore(),
		},

		{
			"Match_NilAndEmptyStringSubject",
			NewNamedNode(""),
			nil,
			nil,
			nil,
			0,
			"Expected Match to return 0 quads when subject is an empty string",
			store,
		},
		{
			"Match_EmptyStringPredicate",
			nil,
			NewNamedNode(""),
			nil,
			nil,
			0,
			"Expected Match to return 0 quads when predicate is an empty string",
			store,
		},
		{
			"Match_EmptyStringObject",
			nil,
			nil,
			NewNamedNode(""),
			nil,
			0,
			"Expected Match to return 0 quads when object is an empty string",
			store,
		},
		{
			"Match_EmptyStringGraph",
			nil,
			nil,
			nil,
			NewNamedNode(""),
			0,
			"Expected Match to return 0 quads when graph is an empty string",
			store,
		},

		{
			"Match_NonExistentTerms",
			NewNamedNode("nonexistent_subject"),
			NewNamedNode("nonexistent_predicate"),
			NewNamedNode("nonexistent_object"),
			NewNamedNode("nonexistent_graph"),
			0,
			"Expected Match to return 0 quads when using terms that do not exist",
			store,
		},
		{
			"Match_SubjectAndNonExistentPredicate",
			NewNamedNode("subject"),
			NewNamedNode("nonexistent_predicate"),
			nil,
			nil,
			0,
			"Expected Match to return 0 quads when subject matches but predicate does not exist",
			store,
		},
		{
			"Match_SubjectAndGraph",
			NewNamedNode("subject"),
			nil,
			nil,
			NewNamedNode("graph"),
			1,
			"Expected Match to return 1 quad when matching subject and graph",
			store,
		},

		{
			"Match_PartialMatchWithNilGraph",
			NewNamedNode("subject1"),
			NewNamedNode("predicate1"),
			nil,
			nil,
			2,
			"Expected Match to return 2 quads when matching subject and predicate, ignoring the graph",
			store,
		},
		{
			"Match_PartialMatchWithEmptyGraph",
			NewNamedNode("subject1"),
			NewNamedNode("predicate1"),
			nil,
			NewNamedNode(""),
			0,
			"Expected Match to return 0 quads when graph is explicitly empty",
			store,
		},
		{
			"Match_AllNilExceptGraph",
			nil,
			nil,
			nil,
			NewNamedNode("graph"),
			5,
			"Expected Match to return 5 quads when matching only graph",
			store,
		},
		{
			"Match_AllNilExceptGraphWithDefault",
			nil,
			nil,
			nil,
			NewDefaultGraph(),
			1,
			"Expected Match to return 1 quad when matching only the default graph",
			store,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := Stream(tt.store.Match(tt.subject, tt.predicate, tt.object, tt.graph)).Count()
			if count != tt.expected {
				t.Errorf("%s, but got %d", tt.message, count)
				for quad := range tt.store.Match(tt.subject, tt.predicate, tt.object, tt.graph) {
					println(quad.ToString())
				}
			}
		})
	}
}

func TestRemoveQuad_MultipleGraphs(t *testing.T) {
	store := NewStore()

	quad1, _ := NewQuad(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph1"),
	)
	quad2, _ := NewQuad(
		NewNamedNode("subject"),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph2"),
	)
	store.AddQuad(quad1)
	store.AddQuad(quad2)

	store.RemoveQuad(quad1)

	if store.Size() != 1 {
		t.Errorf("Expected store size to be 1 after removing one quad, but got %d", store.Size())
	}
	if store.Has(quad1) {
		t.Errorf("Expected quad1 to be removed, but it is still present in the store")
	}
	if !store.Has(quad2) {
		t.Errorf("Expected quad2 to still be in the store, but it is not present")
	}
}

func TestForEach_RemoveDuringIteration(t *testing.T) {
	store := NewStore()

	store.AddQuadFromTerms(
		NewNamedNode("subject1"),
		NewNamedNode("predicate1"),
		NewNamedNode("object1"),
		NewNamedNode("graph1"),
	)
	store.AddQuadFromTerms(
		NewNamedNode("subject2"),
		NewNamedNode("predicate2"),
		NewNamedNode("object2"),
		NewNamedNode("graph2"),
	)

	store.ForEach(func(quad interfaces.IQuad) {
		store.RemoveQuad(quad)
	})

	if store.Size() != 0 {
		t.Errorf(
			"Expected store to be empty after removing all quads during iteration, but got %d",
			store.Size(),
		)
	}
}

func TestImport_NilAndFaultyQuads(t *testing.T) {
	store := NewStore()

	quad1, _ := NewQuad(
		NewNamedNode("subject1"),
		NewNamedNode("predicate1"),
		NewNamedNode("object1"),
		NewNamedNode("graph1"),
	)
	faultyQuad, _ := NewQuad(
		NewLiteral("subject", "", NewNamedNode("")),
		NewNamedNode("predicate"),
		NewNamedNode("object"),
		NewNamedNode("graph"),
	)

	stream := ArrayToStream([]interfaces.IQuad{
		quad1,
		nil,
		faultyQuad,
	})
	store.Import(interfaces.IStream(stream))

	if store.Size() != 1 {
		t.Errorf(
			"Expected store size to be 1 after importing quads with one valid, one nil, and one faulty quad, but got %d",
			store.Size(),
		)
	}
}

func TestDeleteGraph_NonExistentGraph(t *testing.T) {
	store := NewStore()

	quad1, _ := NewQuad(
		NewNamedNode("subject1"),
		NewNamedNode("predicate1"),
		NewNamedNode("object1"),
		NewNamedNode("graph1"),
	)
	store.AddQuad(quad1)

	store.DeleteGraph(NewNamedNode("nonexistent_graph"))

	if store.Size() != 1 {
		t.Errorf(
			"Expected store size to be 1 after deleting a non-existent graph, but got %d",
			store.Size(),
		)
	}
	if !store.Has(quad1) {
		t.Errorf("Expected the quad to still exist in the store, but it was removed")
	}
}

func TestRemoveQuad_LargeStore(t *testing.T) {
	store := NewStore()
	totalQuads := 10_000

	for i := 0; i < totalQuads; i++ {
		store.AddQuadFromTerms(
			NewNamedNode("subject"+strconv.Itoa(i)),
			NewNamedNode("predicate"),
			NewNamedNode("object"),
			NewNamedNode("graph"),
		)
	}

	for i := 0; i < totalQuads; i++ {
		quad, _ := NewQuad(
			NewNamedNode("subject"+strconv.Itoa(i)),
			NewNamedNode("predicate"),
			NewNamedNode("object"),
			NewNamedNode("graph"),
		)
		store.RemoveQuad(quad)
	}

	if store.Size() != 0 {
		t.Errorf("Expected store size to be 0 after removing all quads, but got %d", store.Size())
	}
}

func TestRemoveMatches_LargeStore(t *testing.T) {
	store := NewStore()
	totalQuads := 1_000

	for i := 0; i < totalQuads; i++ {
		store.AddQuadFromTerms(
			NewNamedNode("subject"+strconv.Itoa(i)),
			NewNamedNode("predicate"),
			NewNamedNode("object"),
			NewNamedNode("graph"),
		)
	}

	store.RemoveMatches(nil, nil, nil, nil)

	if store.Size() != 0 {
		t.Errorf("Expected store size to be 0 after removing all quads, but got %d", store.Size())
	}
}

func TestRemove_LargeStore(t *testing.T) {
	store := NewStore()
	totalQuads := 10_000

	for i := 0; i < totalQuads; i++ {
		store.AddQuadFromTerms(
			NewNamedNode("subject"+strconv.Itoa(i)),
			NewNamedNode("predicate"),
			NewNamedNode("object"),
			NewNamedNode("graph"),
		)
	}

	quadStream := make(interfaces.IStream)
	go func() {
		for i := 0; i < totalQuads; i++ {
			quad, _ := NewQuad(
				NewNamedNode("subject"+strconv.Itoa(i)),
				NewNamedNode("predicate"),
				NewNamedNode("object"),
				NewNamedNode("graph"),
			)
			quadStream <- quad
		}
		close(quadStream)
	}()
	store.Remove(quadStream)

	if store.Size() != 0 {
		t.Errorf("Expected store size to be 0 after removing all quads, but got %d", store.Size())
	}
}

func TestConcurrentModifications(t *testing.T) {
	store := NewStore()

	predicate := NewNamedNode("predicate")
	object := NewNamedNode("object")
	graph := NewNamedNode("graph")

	var wg sync.WaitGroup

	// Concurrently add quads
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			store.AddQuadFromTerms(
				NewNamedNode("subject"+strconv.Itoa(i)),
				predicate,
				object,
				graph,
			)
			go func(i int) {
				defer wg.Done()
				quad, _ := NewQuad(
					NewNamedNode("subject"+strconv.Itoa(i)),
					predicate,
					object,
					graph,
				)
				store.RemoveQuad(quad)
			}(i)
		}(i)
	}

	wg.Wait()

	if store.Size() != 0 {
		t.Errorf(
			"Expected store size to be 0 after concurrent modifications, but got %d",
			store.Size(),
		)
	}
}
