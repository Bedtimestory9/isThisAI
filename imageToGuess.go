package main

import (
	"encoding/json"
	"example/user/isThisAI/internal/session"
	"html/template"
	"log"
	"net/http"
	"os"
)

// ImageToGuess type, used as json
type ImageToGuess struct {
	session  *session.Session
	ID       int
	IsThisAI bool
}

func InitImage() *ImageToGuess {
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

	return imageData[i.ID-1] // using traditional index: starting with 0
}

func (i *ImageToGuess) ImageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, i)
}

// AnswerHandler might do too much things right now:
// INFO:
// 1.execute getAnswer() and serialize json
// 2.handle cors with Header().Set()
// 3.execute NextImage()

func (i *ImageToGuess) AnswerHandler(w http.ResponseWriter, r *http.Request) {
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
