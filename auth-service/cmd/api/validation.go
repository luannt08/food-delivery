package api

import (
	"errors"
	"mime"
	"net/http"
	"net/mail"
)

const minPasswordLength = 8

func validateEmailAndPassword(email, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email address")
	}
	if len(password) < minPasswordLength {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

func validateContentType(writer http.ResponseWriter, request *http.Request) bool {
	contenType := request.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contenType)

	if err != nil || mediaType != "application/json" {
		WriteError(writer, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return false
	}
	return true
}
