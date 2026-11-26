package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type MutexImageToGuess struct {
	ID       int
	IsThisAI bool
}

func initImage() *MutexImageToGuess {
	return &MutexImageToGuess{
		ID: 1,
	}
}

func (i *MutexImageToGuess) resetImage() *MutexImageToGuess {
	i.ID = 1
	i.IsThisAI = i.getAnswer().IsThisAI

	return i
}

func (i *MutexImageToGuess) getImage() *MutexImageToGuess {
	i.IsThisAI = i.getAnswer().IsThisAI
	return i
}

// func (i *MutexImageToGuess) nextImage() *MutexImageToGuess {
// 	// If id exceeds maximum image size, it becomes the last image
// 	if i.ID < i.ContainerSize {
// 		i.ID++
// 		i.IsThisAI = i.getAnswer().IsThisAI
// 	}
//
// 	return i
// }

func (i *MutexImageToGuess) getAnswer() *MutexImageToGuess {

	imageData := lookUpAnswerFromFile()
	i.IsThisAI = imageData[i.ID-1].IsThisAI // using traditional index meaning it will start with 0

	return i
}

func lookUpAnswerFromFile() []MutexImageToGuess {
	data, err := os.ReadFile("./dynamic/answer.json")
	if err != nil {
		log.Fatal("Reading from answer json failed!", err)
	}

	var imageData []MutexImageToGuess

	err = json.Unmarshal(data, &imageData)
	if err != nil {
		log.Fatal(err)
	}

	return imageData
}

func (i *MutexImageToGuess) imageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./page/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, i)
}

func (i *MutexImageToGuess) imageAnswerHandler() http.Handler {
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

func (i *MutexImageToGuess) imageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

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
			// i.nextImage()
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

	s := InitSessManager()
	http.Handle("/", s.SessManagerStartSessMiddleware(s.ImageHandler()))
	// http.Handle("/answer", image.imageMiddleware(image.imageAnswerHandler()))
	log.Fatal(http.ListenAndServe(":4321", nil))
}
