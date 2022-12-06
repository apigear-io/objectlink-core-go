package ws

import (
	"github.com/apigear-io/objectlink-core-go/olink/remote"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	registry *remote.Registry
}

func NewServer() *Server {
	return &Server{
		registry: remote.NewRegistry(),
	}
}

func (h *Server) ServeHttp(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn)
	for {
		select {
		case node := <-h.register:
			h.nodes[node] = true
		case node := <-h.unregister:
			if _, ok := h.nodes[node]; ok {
				delete(h.nodes, node)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
