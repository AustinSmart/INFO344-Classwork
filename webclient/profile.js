"use strict"

// ****************** Authorization ******************
window.onload = function() {
    var headers = new Headers();
    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(apiRoot + "users/me",
        {
            method: "GET",
            headers: headers
        })
        .then(function(res){
            if (!res.ok) {
                window.location.href = "index.html";
            } 
            return res.json();
        }).then(function(res) {
            localStorage.setItem("User", JSON.stringify(res));
        })
        .catch(function(err) {
            alert(err);
        });
};
// ****************** End Authorization ******************

// ****************** Variables ******************
var currentUser = JSON.parse(localStorage.getItem("User"));
var userProfile = document.getElementById("user");
    userProfile.innerHTML = currentUser.userName;

var updateForm = document.getElementById("update-form");
var firstNameInput = document.getElementById("first-name-input");
    firstNameInput.value = currentUser.firstName;
var lastNameInput = document.getElementById("last-name-input");
    lastNameInput.value = currentUser.lastName;
var displayNameInput = document.getElementById("display-name-input");
    displayNameInput.value = currentUser.userName;
var emailInput = document.getElementById("email-input");
   emailInput.value = currentUser.email;

// ****************** End Variables ******************

// ****************** Event Listeners ******************
updateForm.addEventListener("submit", function(evt) {
    evt.preventDefault();
    spinner.classList.add("is-active");

    if(validate()) {
        var input = {
                firstName: firstNameInput.value,
                lastName: lastNameInput.value
            };
        var data = JSON.stringify(input)
        var headers = new Headers();
        headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
        fetch(apiRoot + "users/me",
        {
            method: "PATCH",
            body: data,
            headers: headers
        })
        .then(function(res){
            spinner.classList.remove("is-active");
            if (res.ok) {
                alert("Profile updated");
            } else {
                 res.text().then(function (text) {
                    alert(text);
                });
            }   
        }).catch(function(err) {
            spinner.classList.remove("is-active");
            alert(err);
        });
    }
});
// ****************** End Event Listeners ******************

// ****************** Functions ******************
function validate(){
    var valid = true;
    var fields = ["first-name-input", "last-name-input"]
    fields.forEach(function(field) {
        var input = document.forms["update-form"][field].value;
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

function signOut() {
    spinner.classList.add("is-active");
    var headers = new Headers();
    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(apiRoot + "sessions/mine",
        {
            method: "DELETE",
            headers: headers
        })
        .then(function(res){
            if (res.ok) {
                localStorage.clear();
                window.location.href = "index.html";
            } else {
                res.text().then(function (text) {
                    alert("Error: " + text);
                });
            }
        }).catch(function(err) {
            alert(err);
        });
}
// ****************** End Functions ******************
