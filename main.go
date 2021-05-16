package main

import (
	"flag"
	"github.com/rdnply/webrtc-video-conf/app"
	"log"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
)

func main() {
	flag.Parse()

	log.SetFlags(0)

	app := app.New()
	app.RunServer(*addr)
}
