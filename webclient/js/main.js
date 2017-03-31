"use strict";

const API = "http://138.197.219.125/v1/summary?url=";

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
            title.innerHTML = data.title;
            description.innerHTML = data.description;
            image.src = data.image;
            imageDiv.appendChild(image);
            console.log(data)
        }).catch(function (err) {
            console.error(err)
        });
}