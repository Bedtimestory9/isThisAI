package main

import (
	"html/template"
	"log"
	"net/http"
)

// imageToGuess
type imageToGuess struct {
	Id       int
	isThisAI bool
}

func loadImage() *imageToGuess {
	return &imageToGuess{
		Id: 1,
	}
}

func (i *imageToGuess) NextImage() *imageToGuess {
	i.Id++
	return i
}

// template
func renderTemplate(w http.ResponseWriter, i *imageToGuess) {
	t, err := template.ParseFiles("./page/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, i)
}

// page handler
func pageHandler(w http.ResponseWriter, r *http.Request) {
	i := loadImage()
	renderTemplate(w, i)
}

func main() {
	staticFS := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))
	assetFS := http.FileServer(http.Dir("./asset"))
	http.Handle("/asset/", http.StripPrefix("/asset/", assetFS))

	http.HandleFunc("/", pageHandler)
	log.Fatal(http.ListenAndServe(":4321", nil))
}
