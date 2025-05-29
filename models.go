package main

type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type Response struct {
	Message string `json:"message"`
}
