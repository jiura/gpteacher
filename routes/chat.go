package routes

import (
	"html/template"
	"log/slog"
	"net/http"
)

/*
GET /

Loads chat page.
*/
func Chat_GetPage(w http.ResponseWriter, r *http.Request) {
	var err error
	var template *template.Template
	var data struct {
		Messages []struct {
			Author   string
			Nickname string
			Text     string
		}
	}

	// TODO: Get msg history

	template, err = template.ParseFiles("views/_layout.xhtml", "views/chat.xhtml")
	if err != nil {
		goto error
	}

	w.Header().Set("Content-Type", "application/xhtml+xml")
	if r.Header.Get("Hx-Request") == "true" {
		err = template.ExecuteTemplate(w, "body", data)
	} else {
		w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
		err = template.ExecuteTemplate(w, "layout", data)
	}

	if err != nil {
		goto error
	}

	return
error:
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

/*
POST /

Sends message and (partially) reloads chat page (using HTMX).
*/
func Chat_PostMessage(w http.ResponseWriter, r *http.Request) {
	var err error

	var msg string
	var response string

	if err = r.ParseForm(); err != nil {
		goto error
	}

	msg = r.FormValue("message")

	if msg == "" {
		slog.Error("Received message is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: Generate answer and get msg history

	response = `<div class="message bot-message dark-mode"><span class="nickname">Bot</span><br/><pre>Hi! This is a placeholder answer :)</pre></div>`

	w.Header().Set("Content-Type", "application/xhtml+xml")
	w.Write([]byte(response))

	return
error:
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
