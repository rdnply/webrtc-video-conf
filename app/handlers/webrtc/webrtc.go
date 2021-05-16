package webrtc

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/rdnply/webrtc-video-conf/app/rtc"
	"log"
	"net/http"
	"sync"
	"time"
)

type Handle struct {
	connections *rtc.Connections
}

func New(connections *rtc.Connections) *Handle {
	return &Handle{
		connections: connections,
	}
}

func (h *Handle) Handler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Upgrade HTTP request to Websocket
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	c := &rtc.ThreadSafeWriter{unsafeConn, sync.Mutex{}}

	// When this frame returns close the Websocket
	defer c.Close() //nolint

	// Create new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
					"stun:stun2.l.google.com:19302",
					"stun:stun3.l.google.com:19302",
					"stun:stun4.l.google.com:19302",
					"stun:stun.ekiga.net",
					"stun:stun.ideasip.com",
					"stun:stun.schlund.de",
					"stun:stun.stunprotocol.org:3478",
					"stun:stun.voiparound.com",
					"stun:stun.voipbuster.com",
					"stun:stun.voipstunt.com",
				},
			},
			{
				URLs:       []string{"turn:192.158.29.39:3478?transport=udp"},
				Username:   "28224511:1379330808",
				Credential: "JZEOEt2V3Qb0y27GRntt2u2PAYA=",
			},
			{
				URLs:       []string{"turn:192.158.29.39:3478?transport=udp"},
				Username:   "28224511:1379330808",
				Credential: "JZEOEt2V3Qb0y27GRntt2u2PAYA=",
			},
			{
				URLs:       []string{"turn:13.54.1.1:3478?transport=tcp"},
				Username:   "user",
				Credential: "root",
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	})
	if err != nil {
		log.Print(err)
		return
	}

	// When this frame returns close the PeerConnection
	defer peerConnection.Close() //nolint

	// Accept one audio and one video track incoming
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Print(err)
			return
		}
	}

	// Add our new PeerConnection to global list
	h.connections.ListLock.Lock()
	h.connections.PeerConnections = append(h.connections.PeerConnections, rtc.PeerConnectionState{peerConnection, c})
	h.connections.ListLock.Unlock()

	// Trickle ICE. Emit server candidate to client
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}

		if writeErr := c.WriteJSON(&websocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); writeErr != nil {
			log.Println(writeErr)
		}
	})

	// If PeerConnection is closed remove it from global list
	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Println(err)
			}
		case webrtc.PeerConnectionStateClosed:
			h.signalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		// Create a track to fan out our incoming video to all peers
		trackLocal := h.addTrack(t)
		defer h.removeTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})

	peerConnection.OnNegotiationNeeded(func() {
		h.signalPeerConnections()
	})

	// Signal for the new PeerConnection
	h.signalPeerConnections()

	message := &websocketMessage{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}

		switch message.Event {
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Println(err)
				return
			}
		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// Add to list of tracks and fire renegotation for all PeerConnections
func (h *Handle) addTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	h.connections.ListLock.Lock()
	defer func() {
		h.connections.ListLock.Unlock()
		h.signalPeerConnections()
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	h.connections.TrackLocals[t.ID()] = trackLocal
	return trackLocal
}

// Remove from list of tracks and fire renegotation for all PeerConnections
func (h *Handle) removeTrack(t *webrtc.TrackLocalStaticRTP) {
	h.connections.ListLock.Lock()
	defer func() {
		h.connections.ListLock.Unlock()
		h.signalPeerConnections()
	}()

	delete(h.connections.TrackLocals, t.ID())
}

// signalPeerConnections updates each PeerConnection so that it is getting all the expected media tracks
func (h *Handle) signalPeerConnections() {
	h.connections.ListLock.Lock()
	defer func() {
		h.connections.ListLock.Unlock()
		h.connections.DispatchKeyFrame()
	}()

	attemptSync := func() (tryAgain bool) {
		for i := range h.connections.PeerConnections {
			if h.connections.PeerConnections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				h.connections.PeerConnections = append(h.connections.PeerConnections[:i], h.connections.PeerConnections[i+1:]...)
				return true // We modified the slice, start from the beginning
			}

			// map of sender we already are sending, so we don't double send
			existingSenders := map[string]bool{}

			for _, sender := range h.connections.PeerConnections[i].PeerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				// If we have a RTPSender that doesn't map to a existing track remove and signal
				if _, ok := h.connections.TrackLocals[sender.Track().ID()]; !ok {
					if err := h.connections.PeerConnections[i].PeerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range h.connections.PeerConnections[i].PeerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			// Add all track we aren't sending yet to the PeerConnection
			for trackID := range h.connections.TrackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := h.connections.PeerConnections[i].PeerConnection.AddTrack(h.connections.TrackLocals[trackID]); err != nil {
						return true
					}
				}
			}

			offer, err := h.connections.PeerConnections[i].PeerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = h.connections.PeerConnections[i].PeerConnection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			if err = h.connections.PeerConnections[i].Websocket.WriteJSON(&websocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				return true
			}
		}

		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				h.signalPeerConnections()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}
