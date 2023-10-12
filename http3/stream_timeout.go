package http3

import "time"

const StreamTimeout = 60 * time.Minute

type StreamWithDeadline struct {
	Stream
	timeout time.Duration
	lastSet time.Time
}

func NewSteamWithDeadline(stream Stream, timeout time.Duration) Stream {
	return &StreamWithDeadline{
		Stream:  stream,
		timeout: timeout,
		lastSet: time.Time{},
	}
}

func (s *StreamWithDeadline) Read(p []byte) (n int, err error) {
	if time.Since(s.lastSet) >= 10*time.Second {
		_ = s.Stream.SetDeadline(time.Now().Add(s.timeout))
		s.lastSet = time.Now()
	}

	return s.Stream.Read(p)
}

func (s *StreamWithDeadline) Write(p []byte) (n int, err error) {
	if time.Since(s.lastSet) >= 10*time.Second {
		_ = s.Stream.SetDeadline(time.Now().Add(s.timeout))
		s.lastSet = time.Now()
	}

	return s.Stream.Write(p)
}
