package ws

import (
	"fmt"
	"io"
	"sync"
	"time"

	"olink/log"

	"github.com/gorilla/websocket"
)

const (

	// max message size in bytes
	maxMessageSize = 512
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var connId = 0

func nextConnId() string {
	connId++
	return fmt.Sprintf("ws%d", connId)
}

type Connection struct {
	id        string
	socket    *websocket.Conn
	input     chan []byte
	done      chan struct{}
	output    io.WriteCloser
	OnClosing func()
	mu        sync.Mutex
}

func Dial(url string) (*Connection, error) {
	log.Debugf("dial: %s", url)
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	log.Debugf("connected to: %s\n", url)
	conn := NewConnection(ws)
	return conn, nil
}

func NewConnection(socket *websocket.Conn) *Connection {
	p := &Connection{
		id:     nextConnId(),
		socket: socket,
		input:  make(chan []byte),
		done:   make(chan struct{}),
	}
	socket.SetReadLimit(maxMessageSize)
	socket.SetPongHandler(func(string) error {
		deadline := time.Now().Add(pongWait)
		log.Debugf("conn: handle pong %v\n", deadline)
		return socket.SetReadDeadline(deadline)
	})
	socket.SetCloseHandler(func(code int, text string) error {
		// close connection and let write pump handle it
		p.Close()
		return nil
	})
	// SEE /Users/jryannel/work/apigear-go/objectlink-core-go

	go p.ReadPump()
	go p.WritePump()
	return p
}

func (c *Connection) Close() error {
	if c.done == nil {
		return nil
	}
	log.Infof("%s close\n", c.Id())
	if c.output != nil {
		err := c.output.Close()
		if err != nil {
			log.Warnf("%s close output error: %v\n", c.Id(), err)
		}
	}
	log.Debugf("%s: close done\n", c.Id())
	// close done channel to stop write pump
	close(c.done)
	c.done = nil
	c.socket.SetWriteDeadline(time.Now().Add(pongWait))
	err := c.socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Warnf("%s close error: %v\n", c.Id(), err)
	}
	c.socket.Close()
	if c.OnClosing != nil {
		c.OnClosing()
	}
	return nil
}

func (c *Connection) Id() string {
	return c.id
}

func (c *Connection) Url() string {
	return c.socket.RemoteAddr().String()
}

func (c *Connection) SetOutput(out io.WriteCloser) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.output = out
}

func (c *Connection) Write(data []byte) (int, error) {
	log.Debugf("conn: inputC<- %s", data)
	c.input <- data
	return len(data), nil
}

func (c *Connection) WritePump() {
	log.Debugf("%s: start write pump\n", c.Id())
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Debugf("%s: stop write pump\n", c.Id())
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case <-c.done:
			log.Debugf("%s: write pump done\n", c.Id())
			// end go routine
			return
		case data := <-c.input:
			// send message from protocol handler
			log.Debugf("conn: %s <-inputC %s\n", c.Id(), data)
			err := c.socket.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Warnf("write error: %s", err)
				return
			}
		case t := <-ticker.C:
			// send ping message
			log.Debugf("conn: <-ticker %s\n", t)
			err := c.socket.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Warnf("write error: %s", err)
				return
			}
		}
	}
}

func (c *Connection) ReadPump() {
	log.Debugf("%s: start read pump\n", c.Id())
	defer func() {
		log.Debugf("%s: stop read pump\n", c.Id())
		// close connection if we stop reading
		c.Close()
	}()
	for {
		select {
		case <-c.done:
			log.Debugf("conn: <-done\n")
			return
		default:
			c.socket.SetReadDeadline(time.Now().Add(pongWait))
			_, bytes, err := c.socket.ReadMessage()
			if err != nil {
				return
			}
			c.mu.Lock()
			if c.output != nil {
				_, err = c.output.Write(bytes)
				c.mu.Unlock()
			} else {
				log.Warnf("conn: output is nil\n")
				c.mu.Unlock()
				return
			}
			if err != nil {
				log.Warnf("write error: %s", err)
				return
			}
		}
	}
}
