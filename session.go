package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type SessionManager struct {
	l        sync.RWMutex
	sessions map[string]*Session
}

func InitSessManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (s *SessionManager) GetSess(sessid string) (*Session, bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	sess, ok := s.sessions[sessid]
	fmt.Println("goes through getsess, sessions map", s.sessions)
	return sess, ok
}

func (s *SessionManager) SaveSess(sessid string) {
	s.l.Lock()
	defer s.l.Unlock()

	initSess := InitSess(15)
	s.sessions[sessid] = initSess
}

type Session struct {
	SessID        string
	ContainerSize int
	Data          *MutexImageToGuess
}

func InitSess(imageNum int) *Session {
	uuid := uuid.NewString()

	return &Session{
		SessID:        uuid,
		ContainerSize: imageNum,
		Data:          initImage(),
	}
}

type sessIDKey struct{}

func (s *SessionManager) SessIDFromSMContext(ctx context.Context) (string, bool) {
	sessid, ok := ctx.Value(sessIDKey{}).(string)
	return sessid, ok
}

func (s *SessionManager) SessManagerStartSessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessid")

		var cookieAsID string
		if err == nil {
			fmt.Println("error nil goes through here")
			cookieAsID = cookie.String()
			ctx := r.Context()
			ctx = context.WithValue(ctx, sessIDKey{}, cookieAsID)
			r = r.WithContext(ctx)
		} else {
			fmt.Println("error not nil goes through here")
			cookieAsID = uuid.NewString()
			s.SaveSess(cookieAsID)
			ctx := r.Context()
			ctx = context.WithValue(ctx, sessIDKey{}, cookieAsID)
			r = r.WithContext(ctx)
		}

		w.Header().Add("Vary", "Cookie")
		w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)

		next.ServeHTTP(w, r)
	})
}

func (s *SessionManager) ImageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessid, ok := s.SessIDFromSMContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("failed to get session from context")
			return
		}

		session, ok := s.GetSess(sessid)
		if !ok {
			s.SaveSess(sessid)
			session.WriteCookie(w, r)
		}

		t, err := template.ParseFiles("./page/index.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, session)
	})
}

func (s *SessionManager) SessManagerHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Session) WriteCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "sessid",
		Value:    s.SessID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   30_000,
	}

	http.SetCookie(w, cookie)
}

// Session Image Functions

func (s *Session) SessGetImage() *MutexImageToGuess {
	return &MutexImageToGuess{
		ID:       s.Data.ID,
		IsThisAI: s.Data.getAnswer().IsThisAI,
	}
}

func (s *Session) SessNextImage() *MutexImageToGuess {
	// If id exceeds maximum image size, it becomes the last image
	if s.Data.ID < s.ContainerSize {
		s.Data.ID++
		s.Data.IsThisAI = s.Data.getAnswer().IsThisAI
	}
	return s.Data
}

func (s *Session) SessResetImage() *MutexImageToGuess {
	return &MutexImageToGuess{
		ID: initImage().ID,
	}
}

func (s *Session) SessGetImageAnswer() *MutexImageToGuess {
	imageData := lookUpAnswerFromFile()
	s.Data.IsThisAI = imageData[s.Data.ID-1].IsThisAI // using traditional index meaning it will start with 0

	return s.Data
}
