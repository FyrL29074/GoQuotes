package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var quotes map[int]Quote
var nextID int

func main() {
	quotes = make(map[int]Quote)
	nextID = 0

	r := mux.NewRouter()
	r.HandleFunc("/quotes", postQuotes).Methods(http.MethodPost)
	r.HandleFunc("/quotes", getQuotes).Methods(http.MethodGet)
	r.HandleFunc("/quotes/random", getRandomQuote).Methods(http.MethodGet)
	r.HandleFunc("/quotes/{id}", deleteQuoteByID).Methods(http.MethodDelete)
	r.NotFoundHandler = http.HandlerFunc(handle404)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

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

type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type Response struct {
	Message string `json:"message"`
}
