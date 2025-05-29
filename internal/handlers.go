package internal

import (
	"encoding/json"
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

	addQuote(nextID, q)
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
		allQuotes := getAllQuotes()

		json.NewEncoder(w).Encode(allQuotes)
		return
	}

	authorQuotes := getQuotesByAuthor(author)

	json.NewEncoder(w).Encode(authorQuotes)
}

func getRandomQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if getQuotesLength() == 0 {
		http.Error(w, "No quotes available", http.StatusNotFound)
		return
	}

	ids := make([]int, 0, getQuotesLength())
	for id := range getAllQuotes() {
		ids = append(ids, id)
	}

	q := getRandomQuoteFromDb(ids)

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

	if checkIfQuoteInMap(id) {
		deleteQuote(id)
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
