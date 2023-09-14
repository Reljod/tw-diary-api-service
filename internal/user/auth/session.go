package auth

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Reljod/tw-diary-api-service/config"
	"github.com/Reljod/tw-diary-api-service/internal/cache"
	"github.com/google/uuid"
)

type SessionHandler interface {
	CreateNew() (*Session, error)
	IsValid(sessionId string) (bool, error)
	Decode(session string) (*Session, error)
	// Refresh(string) (string, error)
	// Invalidate(string) error
}

type SimpleSessionHandler struct {
	Cache  cache.SessionCache
	Config *config.ConfigSchema
}

type Session struct {
	Id        string `json:"id"`
	Expiry    int64  `json:"maxAge"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"createdAt"`
}

type SessionHandlingError struct{}

func (e *SessionHandlingError) Error() string {
	return "Session internal Error"
}

func (e *SessionHandlingError) Code() string {
	return "50003"
}

func (sessionHandler *SimpleSessionHandler) CreateNew() (*Session, error) {
	// Handles possible duplication of session id
	var retry int32 = 32
	var sessionId string = uuid.New().String()
	for sessionHandler.isSessionIdExists(sessionId) && retry > 0 {
		sessionId = uuid.New().String()
		retry--
	}

	maxAge := sessionHandler.Config.Session.Expiry
	createdAt := time.Now().Unix()
	newSession := &Session{Id: sessionId, Expiry: maxAge, Status: "active", CreatedAt: createdAt}
	newSessionJson, err := json.Marshal(newSession)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil, &SessionHandlingError{}
	}

	err = sessionHandler.Cache.Set(sessionId, string(newSessionJson), nil)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (sessionHandler *SimpleSessionHandler) IsValid(sessionId string) (bool, error) {
	var session Session

	val, err := sessionHandler.Cache.Get(sessionId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return false, &SessionHandlingError{}
	}

	if val == "" {
		return false, nil
	}

	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return false, &SessionHandlingError{}
	}

	expiresAt := time.Unix(session.CreatedAt, 0).UTC().Add(time.Duration(session.Expiry) * time.Millisecond)
	diff := expiresAt.Sub(time.Now().UTC()).Milliseconds()

	if session.Expiry > 0 && session.Status == "active" && diff > 0 {
		return true, nil
	}

	return false, nil
}

func (sessionHandler *SimpleSessionHandler) isSessionIdExists(sessionId string) bool {

	val, err := sessionHandler.Cache.Get(sessionId)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return true
	}

	if val == "" {
		return false
	}

	return true
}

func (sessionHandler *SimpleSessionHandler) Decode(session string) (*Session, error) {
	var cookieSession *Session = nil
	var cookieBytes []byte = make([]byte, len(session))

	n, err := b64.StdEncoding.Decode(cookieBytes, []byte(session))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil, err
	}

	err = json.Unmarshal(cookieBytes[:n], &cookieSession)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil, err
	}

	return cookieSession, nil
}
