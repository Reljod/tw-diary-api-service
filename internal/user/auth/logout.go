package auth

import (
	"fmt"
	"os"
)

type LogoutError struct{}

func (e *LogoutError) Error() string {
	return "Logout error"
}

func (auth *SimpleSessionBasedAuth) Logout(sessionId string) error {

	err := auth.SessionHandler.Invalidate(sessionId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return &LogoutError{}
	}

	return nil
}
