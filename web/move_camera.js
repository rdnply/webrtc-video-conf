var mousePosition;
var offset = [0, 0];
var isDown = false;

let localCamera = document.getElementById('localCamera')

localCamera.addEventListener('mousedown', function (e) {
    isDown = true;
    offset = [
        localCamera.offsetLeft - e.clientX,
        localCamera.offsetTop - e.clientY
    ];
}, true);

document.addEventListener('mouseup', function () {
    isDown = false;
}, true);

document.addEventListener('mousemove', function (event) {
    event.preventDefault();
    if (isDown) {
        mousePosition = {
            x: event.clientX,
            y: event.clientY
        };

        let newLeftPosition = (mousePosition.x + offset[0]);
        let newTopPosition = (mousePosition.y + offset[1]);
        localCamera.style.left = newLeftPosition + 'px';
        localCamera.style.top = newTopPosition + 'px';
        if (newLeftPosition < 0) {
            localCamera.style.left = 0 + 'px';
        }
        if (newTopPosition < 0) {
            localCamera.style.top = 0 + 'px';
        }
        if (newLeftPosition > window.innerWidth - localCamera.offsetWidth) {
            localCamera.style.left = (window.innerWidth - localCamera.offsetWidth) + 'px';
        }
        if (newTopPosition > window.innerHeight - localCamera.offsetHeight) {
            localCamera.style.top = (window.innerHeight - localCamera.offsetHeight) + 'px';
        }
    }
}, true);

