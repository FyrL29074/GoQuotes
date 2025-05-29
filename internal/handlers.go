package internal

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func postQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var q Quote
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	q.ID = nextID

	quotes[nextID] = q
	nextID++

	resp := Response{
		Message: "Quote was successfully saved!",
	}
	json.NewEncoder(w).Encode(resp)
}

func getQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	author := r.URL.Query().Get("author")

	if author == "" {
		allQuotes := make([]Quote, 0)
		for _, q := range quotes {
			allQuotes = append(allQuotes, q)
		}

		json.NewEncoder(w).Encode(allQuotes)
		return
	}

	authorQuotes := make([]Quote, 0)
	for _, q := range quotes {
		if q.Author == author {
			authorQuotes = append(authorQuotes, q)
		}
	}

	json.NewEncoder(w).Encode(authorQuotes)
}

func getRandomQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(quotes) == 0 {
		http.Error(w, "No quotes available", http.StatusNotFound)
		return
	}

	ids := make([]int, 0, len(quotes))
	for id := range quotes {
		ids = append(ids, id)
	}

	randomIndex := rand.Intn(len(ids))
	randomID := ids[randomIndex]

	q := quotes[randomID]
	json.NewEncoder(w).Encode(q)
}

func deleteQuoteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, ok := quotes[id]

	if ok {
		delete(quotes, id)
		resp := Response{
			Message: "Quote was successfully deleted!",
		}

		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := Response{
		Message: "Quote not found",
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(resp)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	resp := Response{
		Message: "Endpoint not found",
	}
	json.NewEncoder(w).Encode(resp)
}
