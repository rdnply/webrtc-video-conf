function processWebRTC(pc, url) {
    let ws = new WebSocket(url)
    pc.onicecandidate = e => {
        if (!e.candidate) {
            return
        }

        ws.send(JSON.stringify({event: 'candidate', data: JSON.stringify(e.candidate)}))
    }

    ws.onclose = function (evt) {
        window.alert("Websocket has closed")
    }

    ws.onmessage = function (evt) {
        let msg = JSON.parse(evt.data)
        if (!msg) {
            return console.log('failed to parse msg')
        }

        switch (msg.event) {
            case 'offer':
                let offer = JSON.parse(msg.data)
                if (!offer) {
                    return console.log('failed to parse offer')
                }
                pc.setRemoteDescription(offer)
                pc.createAnswer().then(answer => {
                    pc.setLocalDescription(answer)
                    ws.send(JSON.stringify({event: 'answer', data: JSON.stringify(answer)}))
                })
                return

            case 'candidate':
                let candidate = JSON.parse(msg.data)
                if (!candidate) {
                    return console.log('failed to parse candidate')
                }

                pc.addIceCandidate(candidate)
        }
    }

    ws.onerror = function (evt) {
        console.log("ERROR: " + evt.data)
    }
}