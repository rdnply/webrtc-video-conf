let isMuted = true;
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


let isEnableVideo = true;
let processVideoBtn = document.getElementById('processVideo')

let silence = () => {
    let ctx = new AudioContext(), oscillator = ctx.createOscillator();
    let dst = oscillator.connect(ctx.createMediaStreamDestination());
    oscillator.start();
    return Object.assign(dst.stream.getAudioTracks()[0], {enabled: false});
}

let black = ({width = 160, height = 120} = {}) => {
    let canvas = Object.assign(document.createElement("canvas"), {width, height});
    canvas.getContext('2d').fillRect(0, 0, width, height);
    let stream = canvas.captureStream();
    return Object.assign(stream.getVideoTracks()[0], {enabled: false});
}

let blackSilence = (...args) => new MediaStream([black(...args), silence()]);

function processVideo() {
    if (isEnableVideo) {
        let video = document.getElementById('localVideo');
        let stream = video.srcObject;

        stream.getVideoTracks().forEach(track => track.stop());

        navigator.mediaDevices.getUserMedia({video: true})
            .then(stream => {
                video.srcObject = blackSilence();

                return Promise.all(pc.getSenders().map(sender =>
                    sender.replaceTrack(blackSilence().getVideoTracks()[0])))
                // .catch(window.alert);
            });
    } else {
        navigator.mediaDevices.getUserMedia({video: true})
            .then(stream => {
                document.getElementById('localVideo').srcObject = stream

                return Promise.all(pc.getSenders().map(sender =>
                    sender.replaceTrack(stream.getVideoTracks()[0])))
                // .catch(window.alert);
            })
    }
    isEnableVideo = !isEnableVideo;
}

processVideoBtn.onclick = processVideo