"use strict";

var signInForm = document.getElementById("signin-form");
var emailInput = document.getElementById("email-input");
var passwordInput = document.getElementById("password-input");
var displayNameInput = document.getElementById("display-name-input");
var spinner = document.getElementById("spinner");

signInForm.addEventListener("submit", function(evt) {
    spinner.classList.add("is-active");
    evt.preventDefault();
     var input = {
                email: emailInput.value,
                password: passwordInput.value,
            };
        var data = JSON.stringify(input)
        fetch(apiRoot + "sessions",
        {
            method: "POST",
            body: data
        })
        .then(function(res){
            spinner.classList.remove("is-active");
            if (res.ok) {
                localStorage.setItem(headerAuthorization, res.headers.get(headerAuthorization));
                 window.location.href = "messages.html";
            } else {
                 res.text().then(function (text) {
                    alert("Error: " + text);
                });
            }   
        }).catch(function(err) {
            spinner.classList.remove("is-active");
            alert(err);
        });
});