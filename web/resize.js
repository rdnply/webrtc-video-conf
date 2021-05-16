function area(increment, count, width, height, margin = 10) {
    let i = w = 0;
    let h = increment * 0.75 + (margin * 2);
    while (i < (count)) {
        if ((w + increment) > width) {
            w = 0;
            h = h + (increment * 0.75) + (margin * 2);
        }
        w = w + increment + (margin * 2);
        i++;
    }
    if (h > height)
        return false;
    else
        return increment;
}

function dish() {
    let margin = 2;
    let remoteVideos = document.getElementById('remoteCameras');
    let width = remoteVideos.offsetWidth - (margin * 2);
    let height = remoteVideos.offsetHeight - (margin * 2);
    let camera = document.getElementsByClassName('remoteCamera');
    let max = 0;

    let i = 1;
    while (i < 5000) {
        let w = area(i, camera.length, width, height, margin);
        if (w === false) {
            max = i - 1;
            break;
        }
        i++;
    }

    max = max - (margin * 2);
    setCameraSize(max, margin);
}

function setCameraSize(width, margin) {
    let cameras = document.getElementsByClassName('remoteCamera');
    for (var s = 0; s < cameras.length; s++) {
        cameras[s].style.width = width + "px";
        cameras[s].style.margin = margin + "px";
        cameras[s].style.height = (width * 0.75) + "px";
    }
}