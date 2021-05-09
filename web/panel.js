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

