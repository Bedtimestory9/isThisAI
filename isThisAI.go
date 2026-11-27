package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
)

type SessionManager struct {
	l        sync.Mutex
	sessions map[string]*Session
}

type Session struct {
	SessID        string
	ContainerSize int
	Data          *ImageToGuess
}

type ImageToGuess struct {
	ID       int
	IsThisAI bool
}

type sessIDKey struct{}

func InitSessManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (s *SessionManager) GetSess(sessid string) (*Session, bool) {
	s.l.Lock()
	defer s.l.Unlock()

	sess, ok := s.sessions[sessid]
	return sess, ok
}

func (s *SessionManager) CreateSess(w http.ResponseWriter, r *http.Request) *Session {
	s.l.Lock()
	defer s.l.Unlock()

	initSess := InitSess(15)
	s.sessions[initSess.SessID] = initSess

	s.WriteCookie(w, r, initSess.SessID)
	return initSess
}

func (s *SessionManager) WriteCookie(w http.ResponseWriter, r *http.Request, sessid string) {
	cookie := &http.Cookie{
		Name:     "sess-cookie",
		Value:    sessid,
		MaxAge:   172_800,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
}

func (s *SessionManager) FillinEmptySess(w http.ResponseWriter, r *http.Request, existedSessid string) {
	s.l.Lock()
	defer s.l.Unlock()

	s.sessions[existedSessid] = &Session{
		ContainerSize: InitSess(15).ContainerSize,
		SessID:        existedSessid,
		Data:          initImage(),
	}
}

func InitSess(imageNum int) *Session {
	uuid := uuid.NewString()

	return &Session{
		SessID:        uuid,
		ContainerSize: imageNum,
		Data:          initImage(),
	}
}

func ContextWithSessID(ctx context.Context, sessid string) context.Context {
	return context.WithValue(ctx, sessIDKey{}, sessid)
}

func SessIDFromContext(ctx context.Context) (string, bool) {
	sessid, ok := ctx.Value(sessIDKey{}).(string)
	return sessid, ok
}

func (s *SessionManager) cookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sess-cookie")
		if err != nil {
			fmt.Println("cookie:", cookie)
			fmt.Println("error goes thru here: ", err)
			session := s.CreateSess(w, r)

			ctx := r.Context()
			ctx = ContextWithSessID(ctx, session.SessID)
			r = r.WithContext(ctx)
		} else {
			cookieAsID := cookie.Value
			fmt.Println("error free goes through here")
			fmt.Println("cookie as ID", cookieAsID)

			ctx := r.Context()
			ctx = ContextWithSessID(ctx, cookieAsID)
			r = r.WithContext(ctx)
		}

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}

func (s *SessionManager) ImageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessid, ok := SessIDFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var session *Session
		session, ok = s.GetSess(sessid)
		if !ok {
			// WARN:
			// This is a temporary fix for in memory data structure
			// it should probably be removed in favor of database
			s.FillinEmptySess(w, r, sessid)
			session, ok = s.GetSess(sessid)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		t, err := template.ParseFiles("./page/index.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, session)
	})
}

// Session Image Functions

func (s *Session) SessGetImage() *ImageToGuess {
	return &ImageToGuess{
		ID:       s.Data.ID,
		IsThisAI: s.Data.getAnswerFromJSON().IsThisAI,
	}
}

func (s *Session) SessNextImage() *Session {
	s.Data.ID++
	s.Data.IsThisAI = s.Data.getAnswerFromJSON().IsThisAI

	return s
}

func (s *Session) SessResetImage() *Session {
	s.Data.ID = initImage().ID
	return s
}

func (s *Session) SessGetImageAnswer() bool {
	imageData := lookUpAnswerFromFile()
	s.Data.IsThisAI = imageData[s.Data.ID-1].IsThisAI // using traditional index meaning it will start with 0

	return s.Data.IsThisAI
}

func initImage() *ImageToGuess {
	return &ImageToGuess{
		ID: 1,
	}
}

func (i *ImageToGuess) getAnswerFromJSON() *ImageToGuess {

	imageData := lookUpAnswerFromFile()
	i.IsThisAI = imageData[i.ID-1].IsThisAI // using traditional index meaning it will start with 0

	return i
}

func lookUpAnswerFromFile() []ImageToGuess {
	data, err := os.ReadFile("./dynamic/answer.json")
	if err != nil {
		log.Fatal("Reading from answer json failed!", err)
	}

	var imageData []ImageToGuess

	err = json.Unmarshal(data, &imageData)
	if err != nil {
		log.Fatal(err)
	}

	return imageData
}

func (s *SessionManager) answerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessid, ok := SessIDFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, ok := s.GetSess(sessid)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		i := session.SessGetImage()

		j, err := json.Marshal(i)
		if err != nil {
			log.Fatal("Marshalling next image failed!", err)
		}
		w.Write(j)

		ctx = ContextWithSessID(ctx, sessid)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (s *SessionManager) afterAnswerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		answerReq, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Reading answer response failed!", err)
		}
		reqObj := &ImageToGuess{}
		err = json.Unmarshal(answerReq, reqObj)
		if err != nil {
			log.Fatal("Unmarshalling answer response to image object failed!", err)
		}

		ctx := r.Context()
		sessid, ok := SessIDFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, ok := s.GetSess(sessid)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		actualAnswer := session.SessGetImageAnswer()

		if reqObj.IsThisAI == actualAnswer {
			if session.Data.ID <= session.ContainerSize {
				session.SessNextImage()
			} else {
				session.SessResetImage()
			}
		} else {
			session.SessResetImage()
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

	http.Handle("/", s.cookieMiddleware(s.ImageHandler()))
	http.Handle("/answer", s.cookieMiddleware(s.answerMiddleware(s.afterAnswerHandler())))
	log.Fatal(http.ListenAndServe(":4321", nil))
}
