package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./page")))
	log.Fatal(http.ListenAndServe(":4321", nil))
}
