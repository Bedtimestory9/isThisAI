package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type MutexImageToGuess struct {
	l             sync.RWMutex
	ID            int
	IsThisAI      bool
	ContainerSize int
}

func initImage(imageNum int) *MutexImageToGuess {
	return &MutexImageToGuess{
		ID:            1,
		ContainerSize: imageNum,
	}
}

func (i *MutexImageToGuess) resetImage() *MutexImageToGuess {
	i.l.Lock()
	defer i.l.Unlock()

	i.ID = 1
	i.IsThisAI = i.appendAnswer().IsThisAI
	return i
}

func (i *MutexImageToGuess) getImage() *MutexImageToGuess {
	i.l.Lock()
	defer i.l.Unlock()

	i.IsThisAI = i.appendAnswer().IsThisAI
	return i
}

func (i *MutexImageToGuess) nextImage() *MutexImageToGuess {
	i.l.Lock()
	defer i.l.Unlock()

	// If id exceeds maximum image size, it becomes the last image
	if i.ID < i.ContainerSize {
		i.ID++
		i.IsThisAI = i.appendAnswer().IsThisAI
	}
	return i
}

func (i *MutexImageToGuess) appendAnswer() *MutexImageToGuess {
	data, err := os.ReadFile("./dynamic/answer.json")
	if err != nil {
		log.Fatal("Reading from answer json failed!", err)
	}

	var imageData []MutexImageToGuess

	err = json.Unmarshal(data, &imageData)
	if err != nil {
		log.Fatal(err)
	}

	i.IsThisAI = imageData[i.ID-1].IsThisAI // using traditional index meaning it will start with 0

	return i
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

func (i *MutexImageToGuess) answerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i.getImage()

		j, err := json.Marshal(i)
		if err != nil {
			log.Fatal("Marshalling next image failed!", err)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(j)
	})
}

func (i *MutexImageToGuess) imageMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		answerReq, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Reading answer response failed!", err)
		}
		reqObj := &MutexImageToGuess{}
		err = json.Unmarshal(answerReq, reqObj)
		if err != nil {
			log.Fatal("Unmarshalling answer response to image object failed!", err)
		}

		actualAnswer := i.getImage().IsThisAI

		if reqObj.IsThisAI == actualAnswer {
			i.nextImage()
		} else {
			i.resetImage()
		}
	})
}

func main() {
	// init static files
	staticFS := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))
	assetFS := http.FileServer(http.Dir("./asset"))
	http.Handle("/asset/", http.StripPrefix("/asset/", assetFS))

	image := initImage(15)

	http.HandleFunc("/", image.imageHandler)
	http.Handle("/answer", image.imageMiddleware(image.answerHandler()))
	log.Fatal(http.ListenAndServe(":4321", nil))
}
