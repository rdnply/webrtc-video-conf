function removeCamera() {
    let cameras = document.getElementsByClassName('remoteCamera');
    let camera = cameras[cameras.length - 1];
    camera.parentNode.removeChild(camera);

    dish();
}

function addCamera(event) {
    let camera = document.createElement('div');
    camera.className = 'remoteCamera';

    let el = document.createElement(event.track.kind)
    el.className = 'remoteVideo'
    el.srcObject = event.streams[0]
    el.autoplay = true
    el.controls = true

    camera.appendChild(el)

    event.track.onmute = function (event) {
        el.play()
    }

    event.streams[0].onremovetrack = ({track}) => {
        if (el.parentNode) {
            el.parentNode.removeChild(el)
            removeCamera()
        }
    }

    let remoteVideos = document.getElementById('remoteCameras');
    remoteVideos.appendChild(camera);

    dish();
}