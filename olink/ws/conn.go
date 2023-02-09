package ws

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
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

var connId atomic.Int32

func nextConnId() string {
	next := connId.Add(1)
	return "c" + strconv.Itoa(int(next))
}

func Dial(ctx context.Context, url string) (*Connection, error) {
	log.Debug().Msgf("dial: %s", url)
	ws, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	log.Debug().Msgf("connected to: %s", url)
	conn := NewConnection(ctx, ws)
	return conn, nil
}

type Connection struct {
	sync.RWMutex
	id            string
	socket        *websocket.Conn
	in            chan []byte
	ctx           context.Context
	ctxCancel     context.CancelFunc
	out           io.WriteCloser
	closeHandlers []func()
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
	c.Lock()
	defer c.Unlock()
	c.closeHandlers = append(c.closeHandlers, onClosing)
}

func (c *Connection) EmitClosing() {
	c.RLock()
	handlers := c.closeHandlers
	c.RUnlock()
	for _, h := range handlers {
		h()
	}
}

func (c *Connection) Close() error {
	c.ctxCancel()
	return nil
}

func (c *Connection) Id() string {
	c.RLock()
	defer c.RUnlock()
	return c.id
}

func (c *Connection) Url() string {
	c.RLock()
	defer c.RUnlock()
	return c.socket.RemoteAddr().String()
}

func (c *Connection) SetOutput(out io.WriteCloser) {
	c.Lock()
	c.out = out
	c.Unlock()
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Close()
		ticker.Stop()
		log.Debug().Msgf("%s: exit write pump ", c.id)
	}()
	for {
		select {
		case <-c.ctx.Done():
			log.Info().Msgf("%s: closing", c.id)
			if c.out != nil {
				c.out.Close()
			}
			c.socket.Close()
			c.EmitClosing()
			return
		case <-ticker.C:
			err := c.socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(sendWait))
			if err != nil {
				log.Error().Msgf("%s: write ping error: %v", c.id, err)
			}
		case bytes := <-c.in:
			log.Debug().Msgf("%s: write: %s", c.id, string(bytes))
			err := c.socket.SetWriteDeadline(time.Now().Add(sendWait))
			if err != nil {
				log.Error().Msgf("%s: set write deadline error: %v", c.id, err)
			}
			err = c.socket.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				log.Error().Msgf("%s: write error: %v", c.id, err)
			}
		}
	}
}

func (c *Connection) ReadPump() {
	defer func() {
		c.Close()
		log.Debug().Msgf("%s: exit read pump ", c.id)
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.socket.SetReadDeadline(time.Now().Add(pongWait))
			_, bytes, err := c.socket.ReadMessage()
			if err != nil {
				log.Info().Msgf("%s: can not read: %v", c.id, err)
				return
			}
			c.RLock()
			out := c.out
			c.RUnlock()
			if out == nil {
				log.Debug().Msgf("%s: no output", c.id)
				continue
			}
			_, err = out.Write(bytes)
			if err != nil {
				log.Debug().Msgf("%s: write error: %v", c.id, err)
			}
		}
	}
}

func (c *Connection) Write(bytes []byte) (int, error) {
	log.Debug().Msgf("%s: write: %s", c.id, string(bytes))
	c.in <- bytes
	return len(bytes), nil
}
