package main

import (
	"log"
	"net/http"
)

func main() {
	// init static files
	staticFS := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))
	assetFS := http.FileServer(http.Dir("./asset"))
	http.Handle("/asset/", http.StripPrefix("/asset/", assetFS))

	image := InitImage()

	http.HandleFunc("/", image.ImageHandler)
	http.HandleFunc("/answer", image.AnswerHandler)
	log.Fatal(http.ListenAndServe(":4321", nil))
}
