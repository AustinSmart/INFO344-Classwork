"use strict";

const API = "http://api.info344.austinsmart.com/v1/summary?url=";

var form = document.getElementById("search-form");
form.addEventListener("submit", function (e) {
    console.log("CLICK");
    e.preventDefault();
    search();
});

function search() {
    var input = document.getElementById("input").value;
    var table = document.getElementById("table");
    var title = document.getElementById("title")
    var description = document.getElementById("description")
    var image = document.createElement("img");
    image.setAttribute("style", "max-width: 200px; max-height: 200px;");
    var imageDiv = document.getElementById("image-div");
    var results = document.getElementById("results");
   
    fetch(API + input)
        .then(function (resp) {
            return resp.json()
        }).then(function (data) {
            results.style.display = "inherit"

            if (data.title)
                title.innerHTML = data.title;
            else
                title.innerHTML = "No OpenGraph title property"

            if (data.description)
                description.innerHTML = data.description;
            else
                description.innerHTML = "No OpenGraph description property"

            if (data.image) {
                image.src = data.image;
                imageDiv.innerHTML = "";
                imageDiv.appendChild(image);
            } else
                imageDiv.innerHTML = "Image not found"

        }).catch(function (err) {
            window.alert(err);
        });
}