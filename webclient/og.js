"use strict";

const API = "https://info344api.austinsmart.com/v1/summary?url=";

// Form submit listener
var form = document.getElementById("search-form");
form.addEventListener("submit", function (e) {
    e.preventDefault();
    search();
});

// Search function for form
function search() {
    var input = document.getElementById("input").value;
    var table = document.getElementById("table");
    var title = document.getElementById("title")
    var description = document.getElementById("description")
    var image = document.createElement("img");
    image.setAttribute("style", "max-width: 200px; max-height: 200px;");
    var imageDiv = document.getElementById("image-div");
    var results = document.getElementById("results");
   
   // Get the requested URL
    fetch(API + input)
        .then(function (resp) {
            return resp.json()
        }).then(function (data) {
            // Make the results div visible
            results.style.display = "inherit"

            // Set results, use placeholder text if not present
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
            // Display errors to the user
            window.alert(err);
        });
}