const toggle = document.querySelector('.toggle');

toggle.addEventListener('change', function() {
    const anon = document.querySelector(".label")
    if (this.checked) {
        anon.textContent = "Anon"
        toggle.value = true
    } else {
        anon.textContent = "Public"
        toggle.value = false
    }
});

let Editor;

$(function () {
    Editor = editormd("editor-md", {
        width: "100%",
        path: '/templates/lib/editormd/lib/',
        saveHTMLToTextarea : true,
        imageUpload : true,
        imageFormats : ["jpg", "jpeg", "gif", "png", "bmp", "webp"],
        imageUploadURL: "/blogs/upload_image",
        onchange: function () {
            $("#output").html("onchange : this.id =>" + this.id + ", markdown =>" + this.getValue());
        }
    });
});

function saveData() {
    const name = document.getElementById("nameModal");
    name.style.display = "flex";
}

function closeNameModal() {
    document.getElementById("nameModal").style.display = "none";
}

function validateTitle() {
    const submitButton = document.getElementById("save-btn");
    const title = document.querySelector(".titleName");
    submitButton.disabled = title.value === "";
}

document.getElementById("save-btn").addEventListener("click", async function (event) {
    event.preventDefault();
    let url, encodedTitle;
    const title = document.querySelector(".titleName");
    if (title.value === "") {
        alert("ERROR: Please enter a title before save!")
        return
    }
    const classification = document.getElementById("dropdown").value
    const isAnonymous = JSON.parse(document.getElementsByClassName("toggle")[0].value)
    console.log(isAnonymous)
    const data = Editor.getMarkdown();
    if (typeof title.defaultValue === "undefined" || title.defaultValue === "") {
        encodedTitle = encodeURIComponent(title.value)
    } else {
        encodedTitle = encodeURIComponent(title.defaultValue)
    }
    url = "/blogs/save/" + encodedTitle
    const currentURL = window.location.href;
    console.log(JSON.stringify({"title": title.value, "content": data, "classification": classification, "isAnonymous": isAnonymous, "preURL": currentURL}))
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({"title": title.value, "content": data, "classification": classification, "isAnonymous": isAnonymous, "preURL": currentURL})
    })
    const responseData = await response.json()
    if (response.status === 200) {
        const title = responseData["response"]
        console.log("/blogs/" + title)
        window.location.href = "/blogs/read/" + title
    } else {
        alert(responseData["response"])
        window.open("/login", "_blank")
    }
})