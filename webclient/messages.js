"use strict";

// ****************** Onload ******************
window.onload = function() {
    // ****************** Authorization ******************
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
            localStorage.setItem("user", JSON.stringify(res));
            var currentUser = JSON.parse(localStorage.getItem("user"));
            var userProfile = document.getElementById("user-profile-link");
                userProfile.innerHTML = currentUser.userName;
        })
        .catch(function(err) {
            alert(err);
        });
        // ****************** End Authorization ******************
    fetch(apiRoot + "users",
        {
            method: "GET",
            headers: headers
        })
        .then(function(res){
            return res.json();
        }).then(function(res) {
            localStorage.setItem("users", JSON.stringify(res));
            users = JSON.parse(localStorage.getItem("users"));
        })
        .catch(function(err) {
            alert(err);
        });
         
    fetch(apiRoot + "channels",
        {
            method: "GET",
            headers: headers
        }).then(function(res){
            return res.json();
        }).then(function(res) {
            localStorage.setItem("channels", JSON.stringify(res));
            channels = JSON.parse(localStorage.getItem("channels"));
            renderSidebar();
        })
        .catch(function(err) {
            alert(err);
        });
};
// ****************** End Onload ******************

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
var headers = new Headers();

var currentChannel = document.getElementById("current-channel");
var currentChannelObj;
var currentChannelID;
var channels;

var messagesContent = document.getElementById("messages-content");
var messages;

var users;

var addMembersDialog = document.getElementById("add-members");
var addMembersOpenButton = document.getElementById("add-members-button");

var addMembersCloseButton = document.getElementById("close-button");
var addMembersSubmitButton = document.getElementById("add-button");

var usersTable = document.getElementById("users-table");

var sidebar = document.getElementById("sidebar");
    
var spinner = document.getElementById("spinner");

var signOutButtonLarge = document.getElementById("sign-out-button-large");
var signOutButtonSmall = document.getElementById("sign-out-button-small");

var inputFormLarge = document.getElementById("message-input-form-large");
var inputFormSmall = document.getElementById("message-input-form-small");

var inputFieldLarge = document.getElementById("message-input-field-large");
var inputFieldSmall = document.getElementById("message-input-field-small");
// ****************** End Variables ******************

// ****************** Event Listeners ******************

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

if (! addMembersDialog.showModal) {
    dialogPolyfill.registerDialog(addMembersDialog);
}

var showDialogHandler = function(event) {
   if(typeof currentChannelID != "undefined") {
        addMembersDialog.showModal();
        populateUsersTable();
   }
};
addMembersCloseButton.addEventListener("click", function() {
    addMembersDialog.close();
});
addMembersSubmitButton.addEventListener("click", function() {
     addMembersDialog.close();
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

function addMessage(body) {
    var message = {
        channelID: currentChannelID,
	    body: body
    }

    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(apiRoot + "messages",
        {
            method: "POST",
            headers: headers,
            body: JSON.stringify(message)
        }).then(function(res){
            return res.json();
        }).then(function(res) {
            messages = JSON.parse(localStorage.getItem(currentChannelID));
            messages.push(res);
            localStorage.setItem(currentChannelID, JSON.stringify(messages));
             messagesContent.innerHTML = "";
            if(messages != null) {
                messages.forEach(function(message) {
                    renderMessage(message);
                });
            }
        })
        .catch(function(err) {
            alert(err);
        });

}

function renderSidebar() {
    if (channels != null) {
        channels.forEach(function(channel) {
            var link = document.createElement("a");
            link.id = channel.id;
            link.classList.add("mdl-navigation__link");
            link.innerHTML = channel.name;
            link.onclick = function() { changeChannel(channel); }
            sidebar.appendChild(link);
        });
    }
}

function changeChannel(channel) {
    currentChannel.innerHTML = channel.name;
    currentChannelID = channel.id;
    currentChannelObj = channel;
    renderChannel(channel.id)
}

function renderChannel(id) {
    messagesContent.innerHTML = "";
    messages = null;
    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(apiRoot + "channels/" + id,
        {
            method: "GET",
            headers: headers
        }).then(function(res){
            return res.json();
        }).then(function(res) {
            localStorage.setItem(id, JSON.stringify(res));
            messages = JSON.parse(localStorage.getItem(id));
            if(messages != null) {
                messages.forEach(function(message) {
                    renderMessage(message);
                });
            }
        })
        .catch(function(err) {
            alert(err);
        });
}

function renderMessage(message) {
    var li = document.createElement("li");
    li.id = message.id;
    li.classList.add("mdl-list__item", "mdl-list__item--three-line");
        var spanPrimary = document.createElement("span");
        spanPrimary.classList.add("mdl-list__item-primary-content");
            var imgAvatar = document.createElement("img");
            imgAvatar.classList.add("material-icons", "mdl-list__item-avatar");
            imgAvatar.src = message.creatorPhotoUrl;
            var spanUserName = document.createElement("span");
            spanUserName.textContent = message.creatorName;
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
   headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(apiRoot + "messages/" + messageID,
        {
            method: "DELETE",
            headers: headers
        }).then(function(res){
            return res.json();
        }).then(function(res) {
            messages = JSON.parse(localStorage.getItem(currentChannelID));
            messages = messages.filter(function(message) {
                return message.id != messageID;
            });
            localStorage.setItem(currentChannelID, JSON.stringify(messages));
            messagesContent.innerHTML = "";
            if(messages != null) {
                messages.forEach(function(message) {
                    renderMessage(message);
                });
            }
        })
        .catch(function(err) {
            alert(err);
        });
} 

function editMessage(messageID) {
    
}

function populateUsersTable() {
    if(users != null) {
        var filteredUsers = users.filter(
            function(user){
                return currentChannelObj.members.indexOf(user.id) < 0;
            });
        usersTable.innerHTML = "";
        filteredUsers.forEach(function(user) {
            renderUser(user)
        });
    }
}

function renderUser(user) {
    var tr = document.createElement("tr");
        var td = document.createElement("td");
            var label = document.createElement("label");
            label.classList.add("mdl-checkbox", "mdl-js-checkbox", "mdl-js-ripple-effect", "mdl-data-table__select");
            label.setAttribute("for", user.id);
                var input = document.createElement("input");
                input.classList.add("mdl-checkbox__input");
                input.setAttribute("type", "checkbox");
                input.id = user.id;
                label.appendChild(input);
            td.appendChild(label);
        tr.appendChild(td);
        td = document.createElement("td");
        td.classList.add("mdl-data-table__cell--non-numeric");
        td.innerHTML = user.userName;
        tr.appendChild(td);
    usersTable.appendChild(tr);
}
// ****************** End Functions ******************