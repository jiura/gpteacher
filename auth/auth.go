package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"

	db "gpteacher/data"
)

func Register(username, password string) error {
	if len(password) > 25 {
		return errors.New("Password for new user is too long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = db.User_Create(username, string(hash))

	return err
}

func Authenticate(username, password, session_token string) error {
	password_hash, err := db.User_ReadPasswordHash(username)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password)); err != nil {
		return err
	}

	return db.UserSession_CreateOrUpdate(username, session_token)
}

func ResetPassword(username, new_password string) error {
	if len(new_password) > 25 {
		return errors.New("New password for user is too long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.User_UpdatePassword(username, string(hash))
}

func CheckSession(cookie *http.Cookie) error {
	username, session_token, found := strings.Cut(cookie.Value, ",")
	if !found {
		return errors.New("Cookie value wrongly formatted")
	}

	return db.UserSession_Check(username, session_token)
}
