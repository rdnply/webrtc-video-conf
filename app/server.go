package app

import (
	"log"
	"net/http"
	"time"
)

func (app *App) RunServer(addr string) {
	// request a keyframe every 3 seconds
	go func() {
		for range time.NewTicker(time.Second * 3).C {
			app.connections.DispatchKeyFrame()
		}
	}()

	log.Fatal(http.ListenAndServeTLS(addr, "server.crt", "server.key", app.routes()))
}
