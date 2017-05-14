package handlers

import (
	"encoding/json"
	"net/http"

	"path"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/messages"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

//ChannelsHandler GET: gets the channels the current user can see
//POST: add the current user to the new channel's Members list, insert the new channel, and returns the newly-inserted Channel
func (ctx *Context) ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	u := users.User{}
	s := SessionState{
		User: &u,
	}
	sid, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &s)
	if err != nil || sid.String() == "" {
		http.Error(w, "Error retrieving state:"+err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		c, err := ctx.MessagesStore.GetAllChannels(u.ID)
		if err != nil {
			http.Error(w, "Error retrieving channels:"+err.Error(), http.StatusInternalServerError)
			return
		}
		respond(w, c)
	case "POST":
		nc := &messages.NewChannel{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(nc)
		if err != nil {
			http.Error(w, "Bad JSON:"+err.Error(), http.StatusBadRequest)
		}
		c, err := ctx.MessagesStore.InsertChannel(u.ID, nc)
		if err != nil {
			http.Error(w, "Error creating channel:"+err.Error(), http.StatusInternalServerError)
			return
		}
		err = ctx.MessagesStore.AddUser(&u.ID, c.ID)
		if err != nil {
			http.Error(w, "Error creating channel:"+err.Error(), http.StatusInternalServerError)
			return
		}
		respond(w, c)
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

	u := users.User{}
	s := SessionState{
		User: &u,
	}
	sid, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &s)
	if err != nil || sid.String() == "" {
		http.Error(w, "Error retrieving state:"+err.Error(), http.StatusInternalServerError)
		return
	}
	_, id := path.Split(r.RequestURI)

	c, err := ctx.MessagesStore.GetChannel(messages.ChannelID(id))
	if err != nil {
		http.Error(w, "Error getting channel:"+err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		isMem := false
		for _, m := range c.Members {
			if m == u.ID {
				isMem = true
			}
		}
		if c.Private || isMem {
			msgs, err := ctx.MessagesStore.GetMessages(500, c.ID)
			if err != nil {
				http.Error(w, "Error getting messages:"+err.Error(), http.StatusInternalServerError)
				return
			}
			respond(w, msgs)
		}
	case "PATCH":
		cu := &messages.ChannelUpdates{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(cu)
		if err != nil {
			http.Error(w, "Bad JSON:"+err.Error(), http.StatusBadRequest)
		}
		if c.CreatorID == u.ID {
			err = ctx.MessagesStore.UpdateChannel(cu, c.ID)
			if err != nil {
				http.Error(w, "Error updating channel:"+err.Error(), http.StatusInternalServerError)
			}
			c.Name = cu.Name
			c.Description = cu.Description
			respond(w, c)
		}
	case "DELETE":
		if c.CreatorID == u.ID {
			err = ctx.MessagesStore.DeleteChannel(c.ID)
			if err != nil {
				http.Error(w, "Error deleting channel:"+err.Error(), http.StatusInternalServerError)
			}
			respond(w, "Channel deleted")
		} else {
			http.Error(w, "You must own channel to delete it", http.StatusUnauthorized)
			return
		}
	case "LINK":
		uu := r.Header.Get("UserToAdd")
		if !c.Private {
			id := users.UserID(uu)
			err = ctx.MessagesStore.AddUser(&id, c.ID)
			if err != nil {
				http.Error(w, "Error adding user:"+err.Error(), http.StatusInternalServerError)
			}
		} else {
			if u.ID == c.CreatorID {
				id := users.UserID(uu)
				err = ctx.MessagesStore.AddUser(&id, c.ID)
				if err != nil {
					http.Error(w, "Error adding user:"+err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.Error(w, "Only private channel owners may add users", http.StatusUnauthorized)
				return
			}
		}
		respond(w, "User has been added")
	case "UNLINK":
		uu := r.Header.Get("UserToRemove")
		if !c.Private {
			id := users.UserID(uu)
			err = ctx.MessagesStore.RemoveUser(&id, c.ID)
			if err != nil {
				http.Error(w, "Error adding user:"+err.Error(), http.StatusInternalServerError)
			}
		} else {
			if u.ID == c.CreatorID {
				id := users.UserID(uu)
				err = ctx.MessagesStore.AddUser(&id, c.ID)
				if err != nil {
					http.Error(w, "Error adding user:"+err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.Error(w, "Only private channel owners may add users", http.StatusUnauthorized)
				return
			}
		}
		respond(w, "User has been removed")
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

	u := users.User{}
	s := SessionState{
		User: &u,
	}
	sid, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &s)
	if err != nil || sid.String() == "" {
		http.Error(w, "Error retrieving state:"+err.Error(), http.StatusInternalServerError)
		return
	}

	nm := &messages.NewMessage{}
	d := json.NewDecoder(r.Body)
	err = d.Decode(nm)
	if err != nil {
		http.Error(w, "Bad JSON:"+err.Error(), http.StatusBadRequest)
	}
	m, err := ctx.MessagesStore.InsertMessage(u, nm)
	if err != nil {
		http.Error(w, "Error creating message:"+err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, m)
}

//SpecificMessageHandler PATCH: if the current user is the message creator, update the specified message and returns the updated Message
//DELETE: if the current user is the message creator, deletes the message
func (ctx *Context) SpecificMessageHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	u := users.User{}
	s := SessionState{
		User: &u,
	}
	sid, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &s)
	if err != nil || sid.String() == "" {
		http.Error(w, "Error retrieving state:"+err.Error(), http.StatusInternalServerError)
		return
	}

	_, id := path.Split(r.RequestURI)
	m, err := ctx.MessagesStore.GetMessage(messages.MessageID(id))
	if err != nil {
		http.Error(w, "Error getting channel:"+err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "PATCH":
		mu := &messages.MessageUpdates{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(mu)
		if err != nil {
			http.Error(w, "Bad JSON:"+err.Error(), http.StatusBadRequest)
		}
		if u.ID == m.CreatorID {
			err = ctx.MessagesStore.UpdateMessage(mu, m.ID)
			if err != nil {
				http.Error(w, "Error updating message:"+err.Error(), http.StatusInternalServerError)
				return
			}
			m.Body = mu.Body
			respond(w, m)
		} else {
			http.Error(w, "Only the message creator may update a message", http.StatusUnauthorized)
			return
		}
	case "DELETE":
		if u.ID == m.CreatorID {
			err = ctx.MessagesStore.DeleteMessage(m.ID)
			if err != nil {
				http.Error(w, "Error deleting message:"+err.Error(), http.StatusInternalServerError)
				return
			}
			respond(w, "Message deleted")
		} else {
			http.Error(w, "Only the message creator may delete a message", http.StatusUnauthorized)
			return
		}
	default:
		http.Error(w, "Request type must be PATCH or DELETE", http.StatusBadRequest)
		return
	}
}
