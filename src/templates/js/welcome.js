function redirect() {
    window.location.href = "/";
}

setTimeout(redirect, 3000);

document.addEventListener("click", redirect);
