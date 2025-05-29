package internal

import "math/rand"

var quotes map[int]Quote
var nextID int

func InitStorage() {
	quotes = make(map[int]Quote)
	nextID = 0
}

func addQuote(id int, q Quote) {
	quotes[nextID] = q
}

func getAllQuotes() []Quote {
	allQuotes := make([]Quote, 0)

	for _, q := range quotes {
		allQuotes = append(allQuotes, q)
	}

	return allQuotes
}

func getQuotesByAuthor(author string) []Quote {
	authorQuotes := make([]Quote, 0)
	for _, q := range quotes {
		if q.Author == author {
			authorQuotes = append(authorQuotes, q)
		}
	}

	return authorQuotes
}

func getRandomQuoteFromDb(ids []int) Quote {
	randomIndex := rand.Intn(len(ids))
	randomID := ids[randomIndex]

	return quotes[randomID]
}

func deleteQuote(id int) {
	delete(quotes, id)
}

func checkIfQuoteInMap(id int) bool {
	_, inMap := quotes[id]
	return inMap
}

func getQuotesLength() int {
	return len(quotes)
}
