package routes

import (
	"html/template"
	"net/http"
)

/*
GET /

Loads chat page.
*/
func Chat_GetPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("views/_layout.xhtml", "views/chat.xhtml")
	if err != nil {
		goto error
	}

	w.Header().Set("Content-Type", "application/xhtml+xml")
	if r.Header.Get("Hx-Request") == "true" {
		err = template.ExecuteTemplate(w, "body", nil)
	} else {
		w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
		err = template.ExecuteTemplate(w, "layout", nil)
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
	//	var err error

	//	msg := r.FormValue("message")

	// TODO: Generate answer
	answer := "Hi! This is a placeholder answer :)"
	// TODO: Get history

	w.Header().Set("Content-Type", "application/xhtml+xml")

	w.Write([]byte(`
		<div class="message bot-message">Bot: ` + answer + `</div>
	`))

	return
	// error:
	//
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
}
