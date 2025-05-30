package internal

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() http.Handler {
	InitStorage()
	return SetupRouter()
}

func TestPostQuote(t *testing.T) {
	r := setupTestServer()

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid quote",
			body:           `{"author" : "Confucius", "quote" : "Life is simple"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing author",
			body:           `{"quote" : "Life is simple"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing quote",
			body:           `{"author" : "Life is simple"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			body:           `{quote : life`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			body:           ``,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(test.body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetQuotes(t *testing.T) {
	r := setupTestServer()

	t.Run("empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var quotes []Quote
		err := json.NewDecoder(rec.Body).Decode(&quotes)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(quotes) != 0 {
			t.Errorf("expected empty list, got %d quotes", len(quotes))
		}
	})

	t.Run("with quotes", func(t *testing.T) {
		quotes := []string{
			`{"author":"Confucius","quote":"Life is simple"}`,
			`{"author":"Aristotle","quote":"Happiness depends on ourselves."}`,
		}
		for _, q := range quotes {
			req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(q))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
		}

		tests := []struct {
			name           string
			url            string
			expectedStatus int
			expectedLength int
		}{
			{
				name:           "get all quotes",
				url:            "/quotes",
				expectedStatus: http.StatusOK,
				expectedLength: 2,
			},
			{
				name:           "get Confucius quotes",
				url:            "/quotes?author=Confucius",
				expectedStatus: http.StatusOK,
				expectedLength: 1,
			},
			{
				name:           "get unknown author",
				url:            "/quotes?author=Socrates",
				expectedStatus: http.StatusOK,
				expectedLength: 0,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, test.url, nil)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)

				if rec.Code != test.expectedStatus {
					t.Errorf("expected status %d, got %d", test.expectedStatus, rec.Code)
				}

				var quotes []Quote
				err := json.NewDecoder(rec.Body).Decode(&quotes)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(quotes) != test.expectedLength {
					t.Errorf("expected %d quotes, got %d", test.expectedLength, len(quotes))
				}
			})
		}
	})
}

func TestGetRandomQuote(t *testing.T) {
	t.Run("empty storage", func(t *testing.T) {
		r := setupTestServer()

		req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})

	t.Run("with multiple quotes", func(t *testing.T) {
		rand.New(rand.NewSource(42))

		r := setupTestServer()

		quotes := []string{
			`{"author":"Confucius","quote":"First quote"}`,
			`{"author":"Plato","quote":"Second quote"}`,
			`{"author":"Socrates","quote":"Third quote"}`,
		}
		for _, q := range quotes {
			req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(q))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
		}

		req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var quote Quote
		err := json.NewDecoder(rec.Body).Decode(&quote)
		if err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if quote.Quote != "Second quote" {
			t.Errorf("expected 'Second quote', got '%s'", quote.Quote)
		}
	})
}

func TestDeleteQuote(t *testing.T) {
	r := setupTestServer()

	body := `{"author":"Confucius","quote":"To be deleted"}`
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{
			name:           "delete existing quote",
			id:             "0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "delete non-existing quote",
			id:             "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid id",
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/quotes/"+test.id, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}
