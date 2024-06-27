package http3

import (
	"github.com/sagernet/quic-go"
	"time"
)

const StreamTimeout = 10 * time.Minute

type StreamWithDeadline struct {
	quic.Stream
	timeout             time.Duration
	manualReadDeadline  time.Time
	manualWriteDeadline time.Time
}

func NewSteamWithDeadline(stream quic.Stream, timeout time.Duration) quic.Stream {
	return &StreamWithDeadline{
		Stream:  stream,
		timeout: timeout,
	}
}

func (s *StreamWithDeadline) Read(p []byte) (n int, err error) {
	if time.Now().After(s.manualReadDeadline) {
		_ = s.Stream.SetDeadline(time.Now().Add(s.timeout))
	}

	return s.Stream.Read(p)
}

func (s *StreamWithDeadline) Write(p []byte) (n int, err error) {
	if time.Now().After(s.manualWriteDeadline) {
		_ = s.Stream.SetDeadline(time.Now().Add(s.timeout))
	}

	return s.Stream.Write(p)
}

func (s *StreamWithDeadline) SetWriteDeadline(t time.Time) error {
	s.manualWriteDeadline = t
	return s.Stream.SetWriteDeadline(t)
}

func (s *StreamWithDeadline) SetReadDeadline(t time.Time) error {
	s.manualReadDeadline = t
	return s.Stream.SetReadDeadline(t)
}

func (s *StreamWithDeadline) SetDeadline(t time.Time) error {
	s.manualWriteDeadline = t
	s.manualReadDeadline = t
	return s.Stream.SetDeadline(t)
}
