package main

import (
	"net/http"
)

func main() {
	InitStorage()

	r := SetupRouter()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
