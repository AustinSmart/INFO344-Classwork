package handlers

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

//Notifier represents a web sockets notifier
type Notifier struct {
	eventq  chan interface{}
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

//NewNotifier constructs a new Notifer.
func NewNotifier() *Notifier {

	n := Notifier{
		eventq:  make(chan interface{}),
		clients: make(map[*websocket.Conn]bool),
		mu:      sync.RWMutex{},
	}
	return &n
}

//Start begins a loop that checks for new events
//and broadcasts them to all web socket clients.
//This function should be called on a new goroutine
//e.g., `go mynotifer.Start()`
func (n *Notifier) Start() {
	for {
		e := <-n.eventq
		n.broadcast(e)
	}
}

//AddClient adds a new web socket client to the Notifer
func (n *Notifier) AddClient(client *websocket.Conn) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	n.clients[client] = true
}

//Notify will add a new event to the event queue
func (n *Notifier) Notify(event interface{}) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	n.eventq <- event
}

//readPump will read all messages (including control messages)
//send by the client and ignore them.
func (n *Notifier) readPump(client *websocket.Conn) {
	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
			break
		}
	}
}

//broadcast sends the event to all client as a JSON-encoded object
func (n *Notifier) broadcast(event interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for c := range n.clients {
		err := c.WriteJSON(event)
		if err != nil {
			c.Close()
			delete(n.clients, c)
		}
	}
}

//WebSocketUpgradeHandler upgrades a request to a websocket
func (ctx *Context) WebSocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {
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

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error establishing connection:"+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Notifier.AddClient(conn)
}
