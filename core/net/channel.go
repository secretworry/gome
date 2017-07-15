package net

import (
	"io"
	"net"
	"sync"
)

type Channel interface {
	ChannelId() uint32
	io.ReadWriteCloser
}

type conn struct {
	conn net.Conn
}

type channelRegistry struct {
	mu sync.Mutex
	channelId uint32
	channels map[uint32] Channel
}

func newChannelRegistry() *channelRegistry {
	return &channelRegistry{channelId: 0, channels: make(map[uint32] Channel)}
}

func (cr *channelRegistry) nextChannelId() uint32 {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.channelId += 1
	return cr.channelId
}

func (cr *channelRegistry) register(ch Channel) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.channels[ch.ChannelId()] = ch
}
