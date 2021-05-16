package rtc

import (
	"github.com/pion/rtcp"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Connections struct {
	ListLock        sync.RWMutex // lock for PeerConnections and TrackLocals
	PeerConnections []PeerConnectionState
	TrackLocals     map[string]*webrtc.TrackLocalStaticRTP
}

type PeerConnectionState struct {
	PeerConnection *webrtc.PeerConnection
	Websocket      *ThreadSafeWriter
}

// Helper to make Gorilla Websockets threadsafe
type ThreadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)

}

func New() *Connections {
	return &Connections{
		TrackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
	}
}

// DispatchKeyFrame sends a keyframe to all PeerConnections, used everytime a new user joins the call
func (c *Connections) DispatchKeyFrame() {
	c.ListLock.Lock()
	defer c.ListLock.Unlock()

	for i := range c.PeerConnections {
		for _, receiver := range c.PeerConnections[i].PeerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = c.PeerConnections[i].PeerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
