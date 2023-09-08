package auth

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type PasswordManager interface {
	Hash(password string) (string, error)
	IsMatch(password string, hashedPassword string) error
}

type BCryptPasswordManager struct{}

func (passwordManager *BCryptPasswordManager) Hash(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return "", err
	}

	return string(hashedPassword), nil
}

func (passwordManager *BCryptPasswordManager) IsMatch(password string, hashedPassword string) error {

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	return nil
}
