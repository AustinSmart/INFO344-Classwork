"use strict";
var signUpForm = document.getElementById("signup-form");
var emailInput = document.getElementById("email-input");
var firstNameInput = document.getElementById("first-name-input");
var lastNameInput = document.getElementById("last-name-input");
var passwordInput = document.getElementById("password-input");
var confirmPasswordInput = document.getElementById("confirm-password-input");
var confirmPasswordField = document.getElementById("confirm-password-field");
var displayNameInput = document.getElementById("display-name-input");
var displayNameField = document.getElementById("display-name-field");
var spinner = document.getElementById("spinner");

passwordInput.addEventListener("keyup", function() {
    confirmPasswordInput.pattern = passwordInput.value;
    if(passwordInput.value !== confirmPasswordInput.value) {
        confirmPasswordField.classList.add("is-invalid");
    } else {
        confirmPasswordField.classList.remove("is-invalid");   
    }
});

confirmPasswordInput.addEventListener("keyup", function() {
    if(passwordInput.value === confirmPasswordInput.value) {
        confirmPasswordField.classList.remove("is-invalid");   
    }
});

signUpForm.addEventListener("submit", function(evt) {
    evt.preventDefault();
    spinner.classList.add("is-active");

    if(validate()) {
        var input = {
                email: emailInput.value,
                password: passwordInput.value,
                passwordconf: confirmPasswordInput.value,
                username: displayNameInput.value,
                firstname: firstNameInput.value,
                lastname: lastNameInput.value
            };
        var data = JSON.stringify(input)
        fetch(httpConnection + apiRoot + "users",
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
    }
});
function validate(){
    var valid = true;
    var fields = ["email-input", "display-name-input", "password-input", "confirm-password-input", "first-name-input", "last-name-input"]
    fields.forEach(function(field) {
        var input = document.forms["signup-form"][field].value;
        if(input == null || input == ""){
            valid = false;
        }
    });
    if(!valid) {
        spinner.classList.remove("is-active");
        alert("Please complete all fields");
        return false;
    } else return true;
}
