"use strict";

// Iterate over all the <a> tags in the navbar, and if that's the active link, give it the "active" class.
const navLinks = document.querySelectorAll("nav a");
for (let i = 0; i < navLinks.length; i++) {
    const link = navLinks[i];
    if (link.getAttribute("href") === window.location.pathname) {
        link.classList.add("active");
        break;
    }
}