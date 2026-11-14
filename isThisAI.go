package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

// ImageToGuess type, used as json
type ImageToGuess struct {
	ID       int
	IsThisAI bool
}

func initImage() *ImageToGuess {
	return &ImageToGuess{
		ID: 1,
	}
}

func (i *ImageToGuess) NextImage() *ImageToGuess {
	i.ID++
	return i
}

func (i *ImageToGuess) getAnswer() ImageToGuess {
	data, err := os.ReadFile("./dynamic/answer.json")
	if err != nil {
		log.Fatal(err)
	}

	var imageData []ImageToGuess
	err = json.Unmarshal(data, &imageData)
	if err != nil {
		log.Fatal(err)
	}

	return imageData[i.ID-1] // using traditional index meaning it will start with 0
}

func (i *ImageToGuess) imageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, i)
}

// INFO: it might do too much things right now:
// 1.execute getAnswer() and serialize json
// 2.handle cors with Header().Set()
// 3.execute NextImage()
func (i *ImageToGuess) answerHandler(w http.ResponseWriter, r *http.Request) {
	a := i.getAnswer()
	j, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	// In terms of cors here, it has to be `localhost`, not `127.0.0.1`
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Write(j)
	i.NextImage()
}

func main() {
	// init static files
	staticFS := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))
	assetFS := http.FileServer(http.Dir("./asset"))
	http.Handle("/asset/", http.StripPrefix("/asset/", assetFS))

	image := initImage()

	http.HandleFunc("/", image.imageHandler)
	http.HandleFunc("/answer", image.answerHandler)
	log.Fatal(http.ListenAndServe(":4321", nil))
}
