package internal

var quotes map[int]Quote
var nextID int

func InitStorage() {
	quotes = make(map[int]Quote)
	nextID = 0
}
