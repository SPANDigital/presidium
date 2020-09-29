/**
 * An image modal that attaches to all images in an article.
 */
export function init() {
    var modal = document.getElementById('presidium-modal');
    var placeholder = document.querySelector('#presidium-modal img');
    var caption = document.querySelector('#presidium-modal .modal-caption');

    window.addEventListener('keyup', function(e) {
        if (e.keyCode === 27) {
            modal.style.display = 'none';
        }
    });

    document.querySelectorAll('#presidium-content article img').forEach(img => {
        img.className += ' ' + 'scalable';
        img.onclick = function() {
            modal.style.display = 'block';
            placeholder.src = this.src;
            caption.innerHTML = this.alt;
        }
    });
}