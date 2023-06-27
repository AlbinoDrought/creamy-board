/**
 * @param {MouseEvent} e
 * @param {HTMLElement} el
 */
function toggleExpand(e, el) {
  e.preventDefault();

  el = el.querySelector('[data-src-expand-to]');
  var src = el.getAttribute('src');
  var nextSrc = el.getAttribute('data-src-expand-to');
  var expandedState = el.getAttribute('data-src-expanded');

  // unless we removeAttribute('src') first,
  // my browser displays the oldSrc until newSrc is fully loaded
  el.removeAttribute('src');
  el.setAttribute('src', nextSrc);
  el.setAttribute('data-src-expand-to', src);
  el.setAttribute('data-src-expanded', expandedState === 'yes' ? 'no' : 'yes');
}

window.addEventListener('load', function () {
  document.querySelectorAll('[data-src-expand-handler]').forEach(function (el) {
    el.addEventListener('click', function (e) {
      toggleExpand(e, el);
    });
  });
});
