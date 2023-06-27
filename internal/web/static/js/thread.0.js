var body = document.querySelector('.form [name="body"]');

/**
 * @param {string} reply 
 */
function citeReply(reply) {
  body.textContent += '>>' + reply + '\n';
}

function performHashQuote() {
  // check for #q1234
  var hashQuotePost = (window.location.hash.match(/q(\d+)$/) || [])[1];
  if (hashQuotePost) {
    citeReply(hashQuotePost);
    body.scrollIntoView();
  }
}

window.addEventListener('hashchange', performHashQuote);
performHashQuote();
