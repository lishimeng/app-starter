package stream

import (
	"bytes"
	"context"
	"io"
)

type PacketProcessor interface {
	Listen(onPacket func(p []byte))
	Data(data []byte) (n int)
}

type Session interface {
	OnData(p []byte, size int, err error)
	Read(p []byte) (size int, err error)
	Write(p []byte) (size int, err error)
	OnIOErr(err func(interface{}))
}

type SessionCtx struct {
	ctx            context.Context
	proxy          io.ReadWriter
	name           string
	baud           int
	packetBuf      *bytes.Buffer
	pp             PacketProcessor
	readSubscriber func([]byte, error)
	reactMode      bool
	onErr          func(interface{})
}

type Option func(s *SessionCtx)

func WithReact() func(s *SessionCtx) {
	return func(s *SessionCtx) {
		s.reactMode = true
	}
}

func WithPacketProcessor(pp PacketProcessor) func(s *SessionCtx) {
	return func(s *SessionCtx) {
		s.pp = pp
	}
}

func NewSession(ctx context.Context, rwc io.ReadWriter, opts ...Option) (s Session) {

	buf := bytes.NewBuffer(nil)
	buf.Grow(1024 * 1024)
	handler := &SessionCtx{
		proxy:     rwc,
		ctx:       ctx,
		packetBuf: buf,
	}
	for _, opt := range opts {
		opt(handler)
	}
	if handler.reactMode {
		go handler.rxLoop()
	}
	s = handler
	return
}

func (s *SessionCtx) OnIOErr(errHandler func(interface{})) {
	s.onErr = errHandler
}
