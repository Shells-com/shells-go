package main

import (
	"net"
	"time"

	"github.com/golang/snappy"
)

type snappyConn struct {
	c net.Conn

	sr *snappy.Reader
	sw *snappy.Writer
}

func (s *snappyConn) Close() error {
	return s.c.Close()
}

func (s *snappyConn) LocalAddr() net.Addr {
	return s.c.LocalAddr()
}

func (s *snappyConn) RemoteAddr() net.Addr {
	return s.c.RemoteAddr()
}

func (s *snappyConn) Read(b []byte) (int, error) {
	if s.sr == nil {
		s.sr = snappy.NewReader(s.c)
	}
	return s.sr.Read(b)
}

func (s *snappyConn) Write(b []byte) (int, error) {
	if s.sw == nil {
		s.sw = snappy.NewWriter(s.c)
	}
	return s.sw.Write(b)
}

func (s *snappyConn) SetDeadline(t time.Time) error {
	return s.c.SetDeadline(t)
}

func (s *snappyConn) SetReadDeadline(t time.Time) error {
	return s.c.SetReadDeadline(t)
}

func (s *snappyConn) SetWriteDeadline(t time.Time) error {
	return s.c.SetWriteDeadline(t)
}
