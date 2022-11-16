package ws

import (
	"net/http"
	"olink/log"
	"olink/pkg/remote"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Hub maintains the set of active peers
// and broadcasts messages to the peers.
type Hub struct {
	// node - source registry
	registry *remote.Registry
	// registered conns
	conns []*Connection
	// inbound messages from the peers
	broadcast chan []byte
	// register new peers
	register chan *Connection
	// unregister peers
	unregister chan *Connection
}

func NewHub(registry *remote.Registry) *Hub {
	h := &Hub{
		registry:   registry,
		broadcast:  make(chan []byte),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		conns:      make([]*Connection, 0),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			log.Info().Msgf("hub: register: %s\n", conn.Id())
			node := remote.NewNode(h.registry)
			node.SetOutput(conn)
			conn.SetOutput(node)
			h.conns = append(h.conns, conn)
		case conn := <-h.unregister:
			log.Info().Msgf("hub: unregister: %s\n", conn.Id())
			for i, c := range h.conns {
				if c == conn {
					h.conns = append(h.conns[:i], h.conns[i+1:]...)
					c.Close()
					break
				}
			}
		case msg := <-h.broadcast:
			log.Info().Msgf("hub: broadcast: %s\n", msg)
			for _, conn := range h.conns {
				select {
				case conn.input <- msg:
				default:
					close(conn.input)
					h.unregister <- conn
				}
			}
		}
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info().Err(err).Msg("error upgrade http call to websocket")
		return
	}
	log.Info().Msgf("new connection: %s\n", socket.RemoteAddr())
	conn := NewConnection(socket)
	conn.OnClosing = func() {
		h.unregister <- conn
	}
	h.register <- conn
}
