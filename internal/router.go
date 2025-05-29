package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/quotes", postQuotes).Methods(http.MethodPost)
	r.HandleFunc("/quotes", getQuotes).Methods(http.MethodGet)
	r.HandleFunc("/quotes/random", getRandomQuote).Methods(http.MethodGet)
	r.HandleFunc("/quotes/{id}", deleteQuoteByID).Methods(http.MethodDelete)

	r.NotFoundHandler = http.HandlerFunc(handle404)
	return r
}
