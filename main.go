package main

import (
	"goquotes/internal"
	"net/http"
)

func main() {
	internal.InitStorage()

	r := internal.SetupRouter()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
