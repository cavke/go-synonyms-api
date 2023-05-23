package service

import (
	"errors"
	"sync"
)

// Synonymer service that is responsible for retrieving and storing synonyms
type Synonymer interface {
	GetSynonym(word string) ([]string, error)
	AddSynonyms(words ...string) error
}

// GraphSynonymer synonym service that holds data in graph
type GraphSynonymer struct {
	mu       sync.RWMutex
	synonyms map[string][]string
}

func NewGraphSynonymer() *GraphSynonymer {
	return &GraphSynonymer{
		synonyms: make(map[string][]string),
	}
}

// GetSynonym returns list of all synonyms for given word
func (s *GraphSynonymer) GetSynonym(word string) ([]string, error) {
	if word == "" {
		return nil, errors.New("empty word param")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]bool)
	s.visit(word, result)

	// Convert the result map to a slice of synonyms
	var synonyms []string
	for synonym := range result {
		// ignore same word
		if word != synonym {
			synonyms = append(synonyms, synonym)
		}
	}

	return synonyms, nil
}

// visit recursively searches synonym structure for given word,
// storing results of already visited words in result map
func (s *GraphSynonymer) visit(word string, result map[string]bool) {
	if result[word] {
		// word was already looked up
		return
	}

	result[word] = true
	for _, synonym := range s.synonyms[word] {
		s.visit(synonym, result)
	}
}

// AddSynonyms adds synonyms to the system
func (s *GraphSynonymer) AddSynonyms(words ...string) error {
	if len(words) == 0 {
		return errors.New("empty words param")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, word := range words {
		if _, ok := s.synonyms[word]; !ok {
			s.synonyms[word] = []string{}
		}

		for _, otherWord := range words {
			if word != otherWord && !containsString(s.synonyms[word], otherWord) {
				s.synonyms[word] = append(s.synonyms[word], otherWord)
			}
		}
	}
	return nil
}

// containsStr returns true if given list contains str
func containsString(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}
