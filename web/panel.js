var isMuted = true;
let processAudioBtn = document.getElementById('processAudio');

function processAudio() {
    let video = document.getElementById('localVideo');
    if (isMuted) {
        video.muted = false;
    } else {
        video.muted = true;
    }
    isMuted = !isMuted;
}

processAudioBtn.onclick = processAudio;


var isEnableVideo = true;
let processVideoBtn = document.getElementById('processVideo')

function processVideo() {
    if (isEnableVideo) {
        let video = document.getElementById('localVideo');
        let stream = video.srcObject;

        stream.getVideoTracks().forEach(track => track.stop());
        video.srcObject = new MediaStream(stream.getAudioTracks());
    } else {
        navigator.mediaDevices.getUserMedia({video: true, audio: true})
            .then(stream => {
                document.getElementById('localVideo').srcObject = stream
                // stream.getTracks().forEach(track => pc.addTrack(track, stream))
            })
    }
    isEnableVideo = !isEnableVideo;
}

processVideoBtn.onclick = processVideo