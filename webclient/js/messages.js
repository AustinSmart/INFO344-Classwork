"use strict";

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
             var currentUser = JSON.parse(localStorage.getItem("User"));
            var userProfile = document.getElementById("user-profile-link");
                userProfile.innerHTML = currentUser.userName;
        })
        .catch(function(err) {
            alert(err);
        });
    fetch(apiRoot + "channels/5912aefcc740d9000105f880",
        {
            method: "GET",
            headers: headers
        }).then(function(res){
            return res.json();
        }).then(function(res) {
            localStorage.setItem("GenMsg", JSON.stringify(res));
            var messages = JSON.parse(localStorage.getItem("GenMsg"));
            renderChannel();
        })
        .catch(function(err) {
            alert(err);
        });
};
// ****************** End Authorization ******************

// ****************** Spinner ******************
function onReady(callback) {
    var intervalID = window.setInterval(checkReady, 1000);
    function checkReady() {
        if (document.getElementsByTagName("body")[0] !== undefined) {
            window.clearInterval(intervalID);
            callback.call(this);
        }
    }
}

function show(id, value) {
    document.getElementById(id).style.display = value ? 'block' : 'none';
}

onReady(function () {
    document.getElementById("page").style.display = "flex";
    spinner.classList.remove("is-active");
});
// ****************** End Spinner ******************

// ****************** Variables ******************
var currentChannel = document.getElementById("current-channel");
    currentChannel.textContent = "#general";
    currentChannel.value = "general";
    
var spinner = document.getElementById("spinner");

var signOutButtonLarge = document.getElementById("sign-out-button-large");
var signOutButtonSmall = document.getElementById("sign-out-button-small");

var inputFormLarge = document.getElementById("message-input-form-large");
var inputFormSmall = document.getElementById("message-input-form-small");

var inputFieldLarge = document.getElementById("message-input-field-large");
var inputFieldSmall = document.getElementById("message-input-field-small");

var messagesContent = document.getElementById("messages-content");

var messages;
// ****************** End Variables ******************

// ****************** Event Listeners ******************
signOutButtonLarge.addEventListener("click", function() {
    signOut();   
});

signOutButtonSmall.addEventListener("click", function() {
    signOut();
});
   
document.getElementById("general-link").addEventListener("click", function() {
    currentChannel.textContent = "#general";
    currentChannel.value = "general"; 
});

inputFormLarge.addEventListener("submit", function(evt) {
    evt.preventDefault();
    addMessage(inputFieldLarge.value);
    inputFieldLarge.value = "";
});

inputFormSmall.addEventListener("submit", function(evt) {
    evt.preventDefault();
    addMessage(inputFieldSmall.value);
    inputFieldSmall.value = "";
});

document.querySelector("main").addEventListener("click", function(event) {
  if (event.target.textContent.toLowerCase() === "delete_forever") {
        deleteMessage(event.target.id);
  } else if (event.target.textContent.toLowerCase() === "mode_edit") {
      editMessage(event.target.id);
  }
});
// ****************** End Event Listeners ******************

// ****************** Functions ******************
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

function addMessage(message) {
    // var message = {
    //     content: message,
    //     createdOn: firebase.database.ServerValue.TIMESTAMP, 
    //     createdBy: {
    //         uid: currentUser.uid,                   
    //         displayName: currentUser.displayName,   
    //         photo: currentUser.photoURL        
    //     }
    // }
    // if(currentChannel.value === "general") messagesGeneralRef.push(message);
}

function renderChannel() {
    messagesContent.innerHTML = "";
    if(messages) {
       messages.forEach(function(message) {
            renderMessage(message);
        });
    }
}

function renderMessage(message) {
    var li = document.createElement("li");
    li.id = message.id;
    li.classList.add("mdl-list__item", "mdl-list__item--three-line");
        var spanPrimary = document.createElement("span");
        spanPrimary.classList.add("mdl-list__item-primary-content");
            var imgAvatar = document.createElement("img");
            imgAvatar.classList.add("material-icons", "mdl-list__item-avatar");
            imgAvatar.src = "";
            var spanUserName = document.createElement("span");
            spanUserName.textContent = message.creatorID;
                var spanTime = document.createElement("span");
                spanTime.classList.add("time");
                spanTime.textContent = " " + message.createdAt;
                spanUserName.appendChild(spanTime);
                if(message.editedAt != "") {
                    var spanEdited = document.createElement("span");
                    spanEdited.classList.add("time");
                    spanEdited.textContent = " (edited " + message.editedAt +")";
                    spanUserName.appendChild(spanEdited);
                }
            var spanBody = document.createElement("span");
            spanBody.classList.add("mdl-list__item-text-body");
            spanBody.textContent = message.body ;
        spanPrimary.appendChild(imgAvatar);
        spanPrimary.appendChild(spanUserName);
        spanPrimary.appendChild(spanBody);
        var spanSecondary = document.createElement("span");
        spanSecondary.classList.add("mdl-list__item-secondary-action");
            var button = document.createElement("button");
            button.classList.add("mdl-button", "mdl-button--icon", "mdl-js-button");
                var iAction = document.createElement("i");
                iAction.classList.add("material-icons");
                iAction.textContent = "delete_forever";
                iAction.id = message.id;
            button.appendChild(iAction);
        spanSecondary.appendChild(button);
            button = document.createElement("button");
            button.classList.add("mdl-button", "mdl-button--icon", "mdl-js-button");
                iAction = document.createElement("i");
                iAction.classList.add("material-icons");
                iAction.textContent = "mode_edit";
                iAction.id = message.id;
            button.appendChild(iAction);
        spanSecondary.appendChild(button);
    li.appendChild(spanPrimary);
    li.appendChild(spanSecondary);

    messagesContent.appendChild(li);    
}

function deleteMessage(messageID) {
    if(currentChannel.value === "general")  {
        // messagesGeneralRef.child(messageID + "/createdBy/displayName").once("value")
        //     .then(function (snapshot) {
        //         if(currentUser.displayName === snapshot.val()) {
        //             if(window.confirm("Delete this message forever?"))
        //                  messagesGeneralRef.child(messageID).remove();
        //         }
        //     }); 
    } 
}

function editMessage(messageID) {
    if(currentChannel.value === "general")  {
        // messagesGeneralRef.child(messageID).once("value")
        // .then(function (snapshot) {
        //     if(currentUser.displayName === snapshot.child("createdBy/displayName").val()) {
        //         var edited = prompt("Edit the message", snapshot.child("content").val());
        //         if(edited) {
        //             messagesGeneralRef.child(messageID).update({"content": edited});
        //             messagesGeneralRef.child(messageID).update({"edited": true});
        //             messagesGeneralRef.child(messageID).update({"editedOn": firebase.database.ServerValue.TIMESTAMP});
        //         }
        //     }
        // });
    } 
}
// ****************** End Functions ******************