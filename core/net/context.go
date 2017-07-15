package net

import (
	"sync"
	"net"
	"context"
	"time"
)

type Conn struct {
	mu sync.Mutex
	c net.Conn
	opts ContextOpts
}

type ContextOpts struct {
	// HeartbeatInterval defines how often keep-alive heartbeat messages should be sent
	// to each connection.
	HeartbeatInterval time.Duration
	// Threshold defines failure detector threshold
	// A low threshold is prone to generate many wrong suspicions but ensures
	// a quick detection in the event of a real crash. Conversely, a high
	// threshold generates fewer mistakes but needs more time to detect actual crashes.
	Threshold float32
}

type Context struct {
	mu sync.Mutex
	conns map[string] Conn
	shutdownCh <-chan interface{}
}

type Server struct {

}

type ServerHandler func(data []byte, channel *Channel)

func NewContext(opts *ContextOpts) *Context {
	return &Context{
		mu: &sync.Mutex{},
		conns: make(map[string] Conn),
		shutdownCh: make(<-chan interface{}),
	}
}


func (c *Context) Shutdown(ctx context.Context) {
	// TODO
}

func (c *Context) Connect(addr string) (channel *Channel, err error) {
	return;
}

func (c *Context) Bind(add string, handler ServerHandler) (server *Server, err error) {
	return;
}
