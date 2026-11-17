// Package session deals with session
package session

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type SessionRequest struct {
	LocalStorageID string
}

type Session struct {
	LocalStorageID [16]byte
	ImageIndex     int
}

func (s *Session) issueSessID() *Session {
	uuid := uuid.New()
	return &Session{
		LocalStorageID: uuid,
		ImageIndex:     1,
	}
}

func (s *Session) getSess(r *http.Request) *Session {
	var sessionObj Session

	enc, err := json.Marshal(r.Body)
	if err != nil {
		log.Fatal("Marshalling request body for session object failed.")
	}

	err = json.Unmarshal(enc, &sessionObj)
	if err != nil {
		log.Fatal("Unmarshalling session object failed.")
	}

	return &Session{
		LocalStorageID: sessionObj.LocalStorageID,
		ImageIndex:     sessionObj.ImageIndex,
	}
}

func (s *Session) HandshakeHandler(w http.ResponseWriter, r *http.Request) {
	s.issueSessID()
	s.getSess(r)
}
