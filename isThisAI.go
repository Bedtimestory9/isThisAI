package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

type MutexImageToGuess struct {
	l        sync.RWMutex
	ID       int
	IsThisAI bool
}

func initImage() *MutexImageToGuess {
	return &MutexImageToGuess{
		ID: 1,
	}
}

func (i *MutexImageToGuess) NextImage() *MutexImageToGuess {
	i.l.Lock()
	defer i.l.Unlock()

	i.ID++
	return i
}

func (i *MutexImageToGuess) getAnswer() *MutexImageToGuess {
	i.l.Lock()
	defer i.l.Unlock()

	data, err := os.ReadFile("./dynamic/answer.json")
	if err != nil {
		log.Fatal(err)
	}

	var imageData []MutexImageToGuess
	err = json.Unmarshal(data, &imageData)
	if err != nil {
		log.Fatal(err)
	}

	return &MutexImageToGuess{
		ID: imageData[i.ID-1].ID, // using traditional index meaning it will start with 0
	}
}

func (i *MutexImageToGuess) imageHandler(w http.ResponseWriter, r *http.Request) {
	i.l.Lock()
	defer i.l.Unlock()

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
func (i *MutexImageToGuess) answerHandler(w http.ResponseWriter, r *http.Request) {
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
