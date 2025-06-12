package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gpteacher/auth"
	db "gpteacher/data"
	"gpteacher/routes"
)

func CheckAuthSession(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("gpteacher-session")
		if err != nil {
			goto unauthorized
		}

		if err = auth.CheckSession(cookie); err != nil {
			goto unauthorized
		}

		next.ServeHTTP(w, r)
		return
	unauthorized:
		var forward_query string
		if r.URL.Path != "" && r.URL.Path != "/" {
			forward_query += r.URL.Path[1:]

			if r.URL.RawQuery != "" {
				forward_query += "?" + r.URL.RawQuery
			}
		}

		redirect_to := "/auth"
		if forward_query != "" {
			redirect_to += "?forward=" + url.QueryEscape(forward_query)
		}

		http.Redirect(w, r, redirect_to, http.StatusSeeOther)
		return
	})
}

func main() {
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// FIXME: Use CheckSession middleware here
	http.HandleFunc("GET /", routes.Chat_GetPage)
	http.HandleFunc("POST /", routes.Chat_PostMessage)

	http.HandleFunc("GET /auth", routes.Auth_GetSignInPage)
	http.HandleFunc("POST /auth", routes.Auth_PostSignIn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// TODO: Log server error here
		}
	}()

	// NOTE: Wait for interrupt signal to gracefully shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop // NOTE: Block until a signal is received

	// NOTE: Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		// TODO: Log shutdown error here
	}

	db.Close()
}
