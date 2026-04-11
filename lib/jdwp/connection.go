package jdwp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

type Conn struct {
	conn net.Conn
}

func Dial(a string) (*Conn, error) {
	log.Printf("start jdwp connection %q", a)
	conn, err := net.Dial("tcp", a)
	if err != nil {
		return nil, fmt.Errorf("invalid connection: %s", err)
	}
	handhake(conn)
	log.Printf("success handhake")
	return &Conn{conn}, nil
}

func (c *Conn) Next() (PackType, error) {
	return nextPack(c.conn)
}

func handhake(rw io.ReadWriter) error {
	msg := []byte("JDWP-Handshake")
	if n, err := rw.Write(msg); n != len(msg) || err != nil {
		return fmt.Errorf("read handhake error count=%d: %s\n", n, err)
	}
	var msgr [14]byte
	r, err := rw.Read(msgr[:])
	if err != nil {
		return fmt.Errorf("handhake errors [read]=%d: %s\n", r, err)
	}
	if !bytes.Equal(msg, msgr[:]) {
		return fmt.Errorf("invalid handhake response")
	}
	return nil
}
