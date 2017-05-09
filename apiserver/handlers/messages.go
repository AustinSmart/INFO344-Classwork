package handlers

import "net/http"

//ChannelsHandler GET: gets the channels the current user can see
//POST: add the current user to the new channel's Members list, insert the new channel, and returns the newly-inserted Channel
func (ctx *Context) ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case "GET":
	case "POST":
	default:
		http.Error(w, "Request type must be GET or POST", http.StatusBadRequest)
		return
	}
}

//SpecificChannelHandler GET: if the user is allowed, returns 500 most recent messages for the channel
//PATCH: if the current user is the channel creator, updates the specified channel's Name/Description and returns the updated Channel
//DELETE: if the current user is the channel creator, delete the specified channel
//LINK:if the specified channel is public, adds the current user to the Members list of the specified channel
//UNLINK: if the specified channel is public, removes the current user from the Members list of the specified channel
func (ctx *Context) SpecificChannelHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case "PATCH":
	case "DELETE":
	case "LINK":
	case "UNLINK":
	default:
		http.Error(w, "Request type must be PATCH, DELETE, LINK, or UNLINK", http.StatusBadRequest)
		return
	}
}

//MessagesHandler POST: inserts the new message, and responds with the newly-inserted Message
func (ctx *Context) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		http.Error(w, "Request type must be POST", http.StatusBadRequest)
		return
	}
}

//SpecificMessageHandler PATCH: if the current user is the message creator, update the specified message and returns the updated Message
//DELETE: if the current user is the message creator, deletes the message
func (ctx *Context) SpecificMessageHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case "PATCH":
	case "DELETE":
	default:
		http.Error(w, "Request type must be PATCH or DELETE", http.StatusBadRequest)
		return
	}
}
