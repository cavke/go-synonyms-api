package service

import (
	"fmt"
	"testing"
)

func TestSynonymInMemoryService_Get_With_Wrong_Input(t *testing.T) {
	service := NewGraphSynonymer()

	// test GetSynonym for an empty input
	result, err := service.GetSynonym("")
	if err == nil {
		t.Errorf("Expected error got nil")
	}
	if result != nil {
		t.Errorf("Expected nil for errored input, got %v", result)
	}
}

func TestSynonymInMemoryService_Add_With_Wrong_Input(t *testing.T) {
	service := NewGraphSynonymer()

	// test AddSynonyms for an empty input
	var emptyInput []string
	err := service.AddSynonyms(emptyInput...)
	if err == nil {
		t.Errorf("Expected error got nil")
	}
}

func TestSynonymInMemoryService_Simple_Add_Get(t *testing.T) {
	service := NewGraphSynonymer()

	// Test GetSynonym for a word that doesn't exist
	result, _ := service.GetSynonym("apple")
	if result != nil {
		t.Errorf("Expected nil for non-existing word, got %v", result)
	}

	// Test AddSynonyms and GetSynonym
	words := []string{"happy", "joyful", "content"}
	_ = service.AddSynonyms(words...)

	words1 := []string{"sad", "bitter", "pessimistic"}
	_ = service.AddSynonyms(words1...)

	fmt.Printf("data: %v\n", service.synonyms)

	result, _ = service.GetSynonym("happy")
	expected := []string{"joyful", "content"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	result, _ = service.GetSynonym("joyful")
	expected = []string{"happy", "content"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	result, _ = service.GetSynonym("content")
	expected = []string{"happy", "joyful"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	result, _ = service.GetSynonym("sad")
	expected = []string{"bitter", "pessimistic"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSynonymInMemoryService_ImmutableInMemoryDB(t *testing.T) {
	service := NewGraphSynonymer()

	// Test AddSynonyms and GetSynonym
	words := []string{"happy", "joyful", "content"}
	_ = service.AddSynonyms(words...)

	result, _ := service.GetSynonym("happy")
	// change first element in result list
	result[0] = "A"
	result, _ = service.GetSynonym("happy")
	expected := []string{"joyful", "content"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

}

func TestSynonymInMemoryService_Transitivity(t *testing.T) {
	service := NewGraphSynonymer()

	// A -> B
	_ = service.AddSynonyms("A", "B")
	fmt.Printf("data: %v\n", service.synonyms)
	// B -> C
	_ = service.AddSynonyms("B", "C")
	fmt.Printf("data: %v\n", service.synonyms)

	// C -> A, B
	result, _ := service.GetSynonym("C")
	fmt.Printf("s(C): %v\n", result)
	expected := []string{"A", "B"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// B -> A, C
	result, _ = service.GetSynonym("B")
	fmt.Printf("s(B): %v\n", result)
	expected = []string{"A", "C"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// A -> B, C
	result, _ = service.GetSynonym("A")
	fmt.Printf("s(A): %v\n", result)
	expected = []string{"B", "C"}
	if !compareStringSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

}

func compareStringSlices(slice1, slice2 []string) bool {
	// Check if the slices are of the same length
	if len(slice1) != len(slice2) {
		return false
	}

	// Create a map to store the element frequencies of slice1
	frequencies := make(map[string]int)

	// Count the occurrences of each element in slice1
	for _, elem := range slice1 {
		frequencies[elem]++
	}

	// Compare the element frequencies of slice1 and slice2
	for _, elem := range slice2 {
		frequencies[elem]--

		// If the frequency becomes negative or the element doesn't exist in slice1, return false
		if frequencies[elem] < 0 {
			return false
		}
	}

	// All elements match and have the same frequencies
	return true
}
