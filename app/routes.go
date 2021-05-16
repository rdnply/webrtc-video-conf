package app

import (
	"github.com/go-chi/chi"
	"github.com/rdnply/webrtc-video-conf/app/handlers/login"
	"github.com/rdnply/webrtc-video-conf/app/handlers/webrtc"
	"log"
	"net/http"
)

func (app *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.HandleFunc("/login", login.New(app.templates.login).Handler)
	r.HandleFunc("/rtc", webrtc.New(app.connections).Handler)

	// index.html handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := app.templates.index.Execute(w, "wss://"+r.Host+"/rtc"); err != nil {
			log.Fatal(err)
		}
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return r
}
