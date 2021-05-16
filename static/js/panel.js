let isMuted = true;
let processAudioBtn = document.getElementById('processAudio');

function controlIconForAudioBtn() {
    let unmutedIcon = document.getElementById("unmutedAudioIcon");
    let mutedIcon = document.getElementById("mutedAudioIcon");
    if (isMuted) {
        unmutedIcon.hidden = "";
        mutedIcon.hidden = "hidden";
    } else {
        unmutedIcon.hidden = "hidden";
        mutedIcon.hidden = "";
    }
}

function processAudio() {
    let video = document.getElementById('localVideo');
    video.muted = !isMuted;
    controlIconForAudioBtn();
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

function controlIconForVideoBtn() {
    let enabledIcon = document.getElementById("enabledVideoIcon");
    let unabledIcon = document.getElementById("unabledVideoIcon");
    if (!isEnableVideo) {
        enabledIcon.hidden = "";
        unabledIcon.hidden = "hidden";
    } else {
        enabledIcon.hidden = "hidden";
        unabledIcon.hidden = "";
    }
}

function processVideo() {
    let video = document.getElementById('localVideo');

    navigator.mediaDevices.getUserMedia({video: true})
        .then(stream => {
            if (!isEnableVideo) {
                stream = blackSilence();
            }

            video.srcObject = stream;
            return Promise.all(pc.getSenders().map(sender =>
                sender.replaceTrack(stream.getVideoTracks()[0])));
        });
    controlIconForVideoBtn();
    isEnableVideo = !isEnableVideo;
}

processVideoBtn.onclick = processVideo