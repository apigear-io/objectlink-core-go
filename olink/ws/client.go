package ws

import (
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/gorilla/websocket"
)

// Client is a websocket client.
type Client struct {
	conn *websocket.Conn
	recv chan core.Message
	send chan core.Message
	done chan struct{}
}

// NewClient creates a new websocket client.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: nil,
		recv: make(chan core.Message),
		send: make(chan core.Message),
		done: make(chan struct{}),
	}
}

// Receiver returns the channel to receive messages.
func (c *Client) Receiver() chan core.Message {
	return c.recv
}

// Connect connects to the server using websocket.
func (c *Client) Connect(url string) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.readPipe()
	go c.writePipe()
	return nil
}

// Close closes the websocket connection.
func (c *Client) Close() error {
	c.done <- struct{}{}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// WriteMessage sends a message to the server.
func (c *Client) WriteMessage(m core.Message) error {
	c.send <- m
	return nil
}

// readPipe reads messages from websocket connection and
// sends them to the receiver channel.
func (c *Client) readPipe() {
	for {
		select {
		case <-c.done:
			return
		default:
			var m core.Message
			err := c.conn.ReadJSON(&m)
			if err != nil {
				return
			}
			c.recv <- m
		}
	}
}

func (c *Client) writePipe() {
	for {
		select {
		case <-c.done:
			return
		case m := <-c.send:
			err := c.conn.WriteJSON(m)
			if err != nil {
				return
			}
		}
	}
}
