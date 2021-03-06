module github.com/rdnply/webrtc-video-conf

// +heroku goVersion go1.12.17
go 1.15

require (
// +heroku install golang.org/x/crypto/ed25519
	github.com/gorilla/websocket v1.4.2
	github.com/pion/rtcp v1.2.6
	github.com/pion/webrtc/v3 v3.0.14
)
