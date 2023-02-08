package ws

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/apigear-io/objectlink-core-go/log"

	"github.com/gorilla/websocket"
)

const (

	// max message size in bytes
	maxMessageSize = 512
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	sendWait   = 3 * time.Second
)

var connId = 0

func nextConnId() string {
	connId++
	return fmt.Sprintf("ws%d", connId)
}

func Dial(ctx context.Context, url string) (*Connection, error) {
	log.Debug().Msgf("dial: %s", url)
	ws, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("connected to: %s", url)
	conn := NewConnection(ctx, ws)
	return conn, nil
}

type Connection struct {
	sync.Mutex
	id        string
	socket    *websocket.Conn
	in        chan []byte
	ctx       context.Context
	ctxCancel context.CancelFunc
	out       io.WriteCloser
	onClosing func()
}

func NewConnection(ctx context.Context, socket *websocket.Conn) *Connection {
	ctx, cancel := context.WithCancel(ctx)
	p := &Connection{
		id:        nextConnId(),
		socket:    socket,
		in:        make(chan []byte),
		ctx:       ctx,
		ctxCancel: cancel,
	}
	socket.SetReadLimit(maxMessageSize)
	socket.SetPongHandler(func(string) error {
		deadline := time.Now().Add(pongWait)
		log.Debug().Msgf("conn: handle pong %v", deadline)
		return socket.SetReadDeadline(deadline)
	})
	socket.SetCloseHandler(func(code int, text string) error {
		p.Close()
		return nil
	})
	go p.WritePump()
	go p.ReadPump()
	return p
}

func (c *Connection) OnClosing(onClosing func()) {
	c.onClosing = onClosing
}

func (c *Connection) EmitClosing() {
	if c.onClosing != nil {
		c.onClosing()
	}
}

func (c *Connection) Close() error {
	c.ctxCancel()
	return nil
}

func (c *Connection) Id() string {
	return c.id
}

// Name returns the name of the connection
func (c *Connection) Name() string {
	return fmt.Sprintf("conn-%s", c.id)
}

func (c *Connection) Url() string {
	return c.socket.RemoteAddr().String()
}

func (c *Connection) SetOutput(out io.WriteCloser) {
	c.Lock()
	defer c.Unlock()
	c.out = out
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			if c.out != nil {
				c.out.Close()
			}
			c.socket.Close()
			c.EmitClosing()
			return
		case <-ticker.C:
			err := c.socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(sendWait))
			if err != nil {
				return
			}
		case bytes := <-c.in:
			log.Debug().Msgf("%s: write: %s", c.Name(), string(bytes))
			err := c.socket.SetWriteDeadline(time.Now().Add(sendWait))
			if err != nil {
				return
			}
			err = c.socket.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				return
			}
		}
	}
}

func (c *Connection) ReadPump() {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.socket.SetReadDeadline(time.Now().Add(pongWait))
			_, bytes, err := c.socket.ReadMessage()
			if err != nil {
				return
			}
			c.Lock()
			out := c.out
			c.Unlock()
			if out != nil {
				_, err = out.Write(bytes)
				if err != nil {
					log.Debug().Msgf("%s: write error: %v", c.Name(), err)
					return
				}
			}
		}
	}
}

func (c *Connection) Write(bytes []byte) (int, error) {
	c.in <- bytes
	return len(bytes), nil
}
