// This script re-maps links so they work on a single page.
// Syntax needs to be ES5 compatible to work with WkHtmltoPDF
var articles = document.querySelectorAll("[data-link]");

function isExternalLink(link) {
  var anchor = document.createElement("a");
  anchor.href = link;
  return window.location.hostname !== anchor.hostname;
}

var links = document.getElementsByTagName("a");

for (var i = 0; i < links.length; i++) {
  var link = links[i];
  var href = link.getAttribute("href");

  if (href == null || href === "") {
    continue;
  }

  if (isExternalLink(href)) {
    continue;
  }

  // Normalize the URL by removing trailing slash before hash
  href = href.replace(baseURL, "").replace(/\/(?=#)/, "");

  var segments = href.split("/").filter(function(x) {
    return x;
  });

  // slugify the segments by replacing slashes and hashes with dashes
  var newAnchor = segments.map(function(x) {
    return x.replace("/", "-").replace("#", "-");
  }).join("-");

  var match = document.getElementById(newAnchor);

  // If the article exists, link to it instead of the original URL
  if (match) {
    link.href = "#" + newAnchor;
  }
  else {
    // If the article doesn't exist, link to the link hash instead
    if (!link.hash == "") {
      link.href = link.hash
    }
  }
}