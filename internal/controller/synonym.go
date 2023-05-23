package controller

import (
	"encoding/json"
	"go-synonyms-api/internal/service"
	"net/http"
)

type SynonymController struct {
	service.Synonymer
}

type createSynonymRequest struct {
	Synonyms []string
}

// GetSynonym returns synonyms for given word
func (c *SynonymController) GetSynonym(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	word := req.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "word param is empty", http.StatusBadRequest)
		return
	}

	synonyms, err := c.Synonymer.GetSynonym(word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if synonyms == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(synonyms); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateSynonyms adds synonyms to the system.
func (c *SynonymController) CreateSynonyms(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	createSynReq := createSynonymRequest{}
	err := json.NewDecoder(req.Body).Decode(&createSynReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate input
	if len(createSynReq.Synonyms) == 0 {
		http.Error(w, "please provide synonyms list", http.StatusBadRequest)
		return
	} else if len(createSynReq.Synonyms) > 10 {
		http.Error(w, "max 10 synonyms allowed in one request", http.StatusBadRequest)
		return
	}

	err = c.AddSynonyms(createSynReq.Synonyms...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
