package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/assets/upbit-sg/201911/", http.StripPrefix("/assets/upbit-sg/201911", fs))
	err := http.ListenAndServe(":2027", nil)
	if err != nil {
		log.Fatal(err)
	}
}
