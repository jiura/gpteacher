package routes

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"time"

	"gpteacher/auth"
)

func Auth_GetSignInPage(w http.ResponseWriter, r *http.Request) {
	var template *template.Template
	var buf bytes.Buffer

	template, err := template.ParseFiles("views/_layout.xhtml", "views/auth.xhtml")
	if err != nil {
		goto error
	}

	if err = template.ExecuteTemplate(&buf, "layout", nil); err != nil {
		goto error
	}

	w.Header().Set("Content-Type", "application/xhtml+xml")
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	if err = template.ExecuteTemplate(w, "layout", nil); err != nil {
		goto error
	}

	return
error:
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func Auth_PostSignIn(w http.ResponseWriter, r *http.Request) {
	var err error

	var random_bytes [32]byte
	var cookie http.Cookie

	username := r.FormValue("username")
	password := r.FormValue("password")

	_, err = rand.Read(random_bytes[:])
	session_token := base64.URLEncoding.EncodeToString(random_bytes[:])

	if err = auth.Authenticate(username, password, session_token); err != nil {
		goto error
	}

	cookie = http.Cookie{
		Name:     "gpteacher-session",
		Value:    username + "," + session_token,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "/")

	return
error:
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
