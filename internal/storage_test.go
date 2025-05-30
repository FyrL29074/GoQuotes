package internal

import (
	"math/rand"
	"testing"
)

func setup() {
	InitStorage()
}

func TestAddQuote(t *testing.T) {
	setup()

	q := Quote{Author: "Confucius", Quote: "Life is simple"}
	addQuote(q)

	quotes := getAllQuotes()
	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	if quotes[0].Author != q.Author || quotes[0].Quote != q.Quote {
		t.Errorf("stored quote mismatch: got %+v", quotes[0])
	}
}

func TestGetCurrentId(t *testing.T) {
	setup()

	if id := getCurrentId(); id != 0 {
		t.Errorf("expected ID 0, got %d", id)
	}

	addQuote(Quote{Author: "A", Quote: "One"})
	addQuote(Quote{Author: "B", Quote: "Two"})
	addQuote(Quote{Author: "C", Quote: "C"})
	if id := getCurrentId(); id != 3 {
		t.Errorf("expected ID 1, got %d", id)
	}
}

func TestGetQuotesByAuthor(t *testing.T) {
	setup()

	addQuote(Quote{Author: "A", Quote: "One"})
	addQuote(Quote{Author: "B", Quote: "Two"})
	addQuote(Quote{Author: "A", Quote: "Three"})

	quotesByAuthor := getQuotesByAuthor("A")
	if len(quotesByAuthor) != 2 {
		t.Fatalf("expected 2 quotes for author A, got %d", len(quotesByAuthor))
	}
}

func TestGetRandomQuoteFromDb(t *testing.T) {
	setup()
	rand.Seed(42)

	addQuote(Quote{Author: "A", Quote: "One"})
	addQuote(Quote{Author: "B", Quote: "Two"})
	addQuote(Quote{Author: "C", Quote: "Three"})

	ids := []int{0, 1, 2}
	q := getRandomQuoteFromDb(ids)

	if q.Quote != "Three" {
		t.Errorf("expected quote 'Three', got '%s'", q.Quote)
	}
}

func TestDeleteQuoteFromDb(t *testing.T) {
	setup()

	addQuote(Quote{Author: "A", Quote: "To delete"})

	deleteQuote(0)
	if checkIfQuoteInMap(0) {
		t.Error("quote should have been deleted")
	}
}

func TestGetQuotesLength(t *testing.T) {
	setup()
	if getQuotesLength() != 0 {
		t.Error("expected length 0")
	}
	addQuote(Quote{Author: "A", Quote: "One"})
	addQuote(Quote{Author: "B", Quote: "Two"})

	if l := getQuotesLength(); l != 2 {
		t.Errorf("expected length 2, got %d", l)
	}
}
