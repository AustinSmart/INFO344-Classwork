"use strict";

// ****************** Onload ******************
window.onload = function() {
    // ****************** Authorization ******************
    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(httpConnection + apiRoot + "users/me",
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
    fetch(httpConnection + apiRoot + "users",
        {
            method: "GET",
            headers: headers
        })
        .then(function(res){
            if(res.ok) {
                return res.json();
            } else {
                res.text().then(function (text) {
                    alert("Error: " + text);
                });
            }
        }).then(function(res) {
            localStorage.setItem("users", JSON.stringify(res));
            users = JSON.parse(localStorage.getItem("users"));
        })
        .catch(function(err) {
            alert(err);
        });
         
    fetch(httpConnection + apiRoot + "channels",
        {
            method: "GET",
            headers: headers
        }).then(function(res){
            if(res.ok) {
                return res.json();
            } else {
                res.text().then(function (text) {
                    alert("Error: " + text);
                });
            }
        }).then(function(res) {
            localStorage.setItem("channels", JSON.stringify(res));
            channels = JSON.parse(localStorage.getItem("channels"));
            renderSidebar();
        })
        .catch(function(err) {
            alert(err);
        });

    var websocket = new WebSocket(webSocketConnection + apiRoot+ "websocket?" + headerAuthorization + "=" + localStorage.getItem(headerAuthorization));
    websocket.onmessage =  function(wsevent) {
        var msg = JSON.parse(event.data);
        switch(msg.type) {
            case "New Message":
            case "Message Updated":
                var currChannel = channels.filter(function(channel) {
                        return channel.id == msg.Message.channelID;
                    });
                if(currChannel.length > 0) {
                    var saveMessages = JSON.parse(localStorage.getItem(msg.Message.channelID));
                    saveMessages.push(msg.Message);
                    localStorage.setItem(msg.Message.channelID, JSON.stringify(saveMessages));
                    messages = JSON.parse(localStorage.getItem(currentChannelID));
                    if(messages != null) {
                        messagesContent.innerHTML = "";
                        messages.forEach(function(message) {
                            renderMessage(message);
                        });
                    }
                }
                break;
            case "Message Deleted":
                var currChannel = channels.filter(function(channel) {
                        return channel.id == msg.Message.channelID;
                    });
                if(currChannel.length > 0) {
                    messages = JSON.parse(localStorage.getItem(currentChannelID));
                    messages = messages.filter(function(message) {
                        return message.id != msg.Message.id;
                    });
                    if(messages != null) {
                        messagesContent.innerHTML = "";
                        messages.forEach(function(message) {
                            renderMessage(message);
                        });
                    }
                }
                break;
            case "New Channel":
                //TODO add private channel member check
                if(msg.channel.private === false ) {
                    channels = JSON.parse(localStorage.getItem("channels"));
                    channels.push(msg.channel);
                    localStorage.setItem("channels", JSON.stringify(channels));
                    renderSidebar();
                }
                break;
            case "Channel Updated":
                if(msg.channel.private === false ) {
                    channels = JSON.parse(localStorage.getItem("channels"));
                    var filteredChannels = channels.filter(
                        function(channel){
                            return channel.id !== msg.channel.id;
                    });
                    filteredChannels.push(msg.channel);
                    localStorage.setItem("channels", JSON.stringify(filteredChannels));
                    renderSidebar();
                }
                break;
            case "Channel Deleted":
                if(msg.channel.private === false ) {
                    channels = JSON.parse(localStorage.getItem("channels"));
                    var filteredChannels = channels.filter(
                        function(channel){
                            return channel.id !== msg.channel.id;
                    });
                    localStorage.setItem("channels", JSON.stringify(filteredChannels));
                    renderSidebar();
                }
                break;
            case "New User":
                break;
            case "User Joined Channel":
                break;
            case "User Left Channel":
                break;
        };
    };
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
// var currentChannelDescription = document.getElementById("current-channel-description");
var currentChannelObj;
var currentChannelID;
var channels;

var messagesContent = document.getElementById("messages-content");
var messages;

var users;
var addUsers;

var addNewChannelDialog = document.getElementById("add-new-channel");
var addChannelDialog = document.getElementById("add-channel");

var channelCloseButton = document.getElementById("channel-close-button");
var closeButton = document.getElementById("close-button");

var addChannelSubmitButton = document.getElementById("add-channel-button");
var addNewChannelSubmitButton = document.getElementById("add-button");

var usersTableBody = document.getElementById("users-table");
var usersTableMain = document.getElementById("users-table-main");
var membersTable;

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

if (!addChannelDialog.showModal) {
    dialogPolyfill.registerDialog(addChannelDialog);
}

var showChannelDialogHandler = function(event) {
   addChannelDialog.showModal();
};

if (!addNewChannelDialog.showModal) {
    dialogPolyfill.registerDialog(addNewChannelDialog);
}

var showNewChannelDialogHandler = function(event) {
    addNewChannelDialog.showModal();
};

closeButton.addEventListener("click", function() {
    addNewChannelDialog.close();
});

channelCloseButton.addEventListener("click", function() {
    addChannelDialog.close();
});

addNewChannelSubmitButton.addEventListener("click", function() {
    var name = document.getElementById("new-channel-form-name");
    var description = document.getElementById("new-channel-form-description");
    var newChannel = {
        Name: name.value,
        Description: description.value
    }
    createChannel(newChannel);
    addNewChannelDialog.close();
});

addChannelSubmitButton.addEventListener("click", function() {
     addChannelDialog.close();
});
// ****************** End Event Listeners ******************

// ****************** Functions ******************
function signOut() {
    spinner.classList.add("is-active");
    var headers = new Headers();
    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(httpConnection + apiRoot + "sessions/mine",
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

    if(!body.startsWith("bot:")) {
        var message = {
            channelID: currentChannelID,
            body: body
        }

        headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
        fetch(httpConnection + apiRoot + "messages",
            {
                method: "POST",
                headers: headers,
                body: JSON.stringify(message)
            }).then(function(res){
                if(res.ok) {
                    return res.json();
                } else {
                    res.text().then(function (text) {
                        alert("Error: " + text);
                    });
                }
            }).catch(function(err) {
                alert(err);
            });
    } else {
        //CHATBOT
        headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
        fetch(httpConnection + apiRoot + "bot",
            {
                method: "POST",
                headers: headers,
                body: body
            }).then(function(res){
                if(res.ok) {
                    return res.text();
                } else {
                    res.text().then(function (text) {
                        alert("Error: " + text);
                    });
                }
            }).then(function(res) {
                alert(res);
            }).catch(function(err) {
                alert(err);
            });
    }
}

function renderSidebar() {
    var profileLink = sidebar.children[0].cloneNode(true);
    sidebar.innerHTML = "";
    sidebar.appendChild(profileLink);
    channels = JSON.parse(localStorage.getItem("channels"));
    if (channels != null) {
        channels.forEach(function(channel) {
            var link = document.createElement("a");
            link.id = channel.id;
            link.classList.add("mdl-navigation__link");
            link.innerHTML = channel.name;
            link.onclick = function() { changeChannel(channel); }
            sidebar.appendChild(link);
        });
    // if(users != null && typeof currentChannelID != "undefined") { 
    //     var table = document.createElement("table");
    //     table.classList.add("mdl-data-table");
    //         var thead = document.createElement("thead");
    //             var tr = document.createElement("tr");
    //                 var th = document.createElement("th");
    //                 th.innerHTML = "Members";
    //                 tr.appendChild(th);
    //             thead.appendChild(tr);
    //         table.appendChild(thead);
    //             var tbody = document.createElement("tbody");
    //             tbody.id = "members-table";
    //         table.appendChild(tbody);
    //     sidebar.appendChild(table);
    //     membersTable = document.getElementById("members-table");
    //     currentChannelObj.members.forEach(function(member) {
    //         renderMember(member)
    //     });
    // }
    
    }
}

function changeChannel(channel) {
    currentChannel.innerHTML = channel.name;
    // currentChannelDescription.innerHTML = channel.description;
    currentChannelID = channel.id;
    currentChannelObj = channel;
    renderChannel(channel.id)
    // renderSidebar();
}

function renderChannel(id) {
    messagesContent.innerHTML = "";
    messages = null;
    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(httpConnection + apiRoot + "channels/" + id,
        {
            method: "GET",
            headers: headers
        }).then(function(res){
           if(res.ok) {
                return res.json();
            } else {
                res.text().then(function (text) {
                    alert("Error: " + text);
                });
            }
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
    fetch(httpConnection + apiRoot + "messages/" + messageID,
        {
            method: "DELETE",
            headers: headers
        }).then(function(res){
           if(res.ok) {
                return res.json();
            } else {
                res.text().then(function (text) {
                    alert("Error: " + text);
                });
            }
        }).catch(function(err) {
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
        usersTableBody.innerHTML = "";
        filteredUsers.forEach(function(user) {
            renderUser(user)
        });
        addUsers = [];
        var checkboxes = usersTableBody.querySelectorAll('.mdl-checkbox__input');
        for (var i = 0; i < checkboxes.length; i++) {
            checkboxes[i].addEventListener('change', function(e) {
                if (!e.target.tagName === 'input' || e.target.getAttribute('type') !== 'checkbox') { 
                    return;
                }
                addUsers.push(this.id);
            });
        }
    }
}

function renderUser(user) {
    var tr = document.createElement("tr");
    tr.id = user.id;
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
    usersTableBody.appendChild(tr);
}

// function renderMember(member) {
//      var tr = document.createElement("tr");
//         var td = document.createElement("td");
//         td.classList.add("mdl-data-table__cell--non-numeric");
//         var filteredUsers = users.filter(function(user) {
//                 return user.id == member;
//             });
//         td.innerHTML = filteredUsers[0].userName;
//         tr.appendChild(td);
//     membersTable.appendChild(tr);
// }

function createChannel(newChannel) {
    headers.set(headerAuthorization, localStorage.getItem(headerAuthorization));
    fetch(httpConnection + apiRoot + "channels",
        {
            method: "POST",
            headers: headers
        }).then(function(res){
           if(res.ok) {
                return res.json();
            } else {
                res.text().then(function (text) {
                    alert("Error: " + text);
                });
            }
        }).then(function(res) {
            messages = JSON.parse(localStorage.getItem(currentChannelID));
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
// ****************** End Functions ******************