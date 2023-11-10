package stream

import (
	"bytes"
	"context"
	"github.com/tarm/serial"
	"io"
	"time"
)

type Session interface {
	OnData(p []byte, size int, err error)
	Read(p []byte) (size int, err error)
	Write(p []byte) (size int, err error)
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

func NewSerialSession(ctx context.Context, name string, baud int, opts ...Option) (s Session, err error) {
	c := &serial.Config{
		Name:        name,
		Baud:        baud,
		Size:        8,
		ReadTimeout: time.Second * 5,
		StopBits:    1,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return
	}
	s = NewSession(ctx, port, opts...)
	return
}

func NewSession(ctx context.Context, rwc io.ReadWriter, opts ...Option) (s Session) {

	buf := bytes.NewBuffer(nil)
	buf.Grow(1024)
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

func (s *SessionCtx) Close() {

}
