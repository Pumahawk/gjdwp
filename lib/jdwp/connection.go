package jdwp

import (
	"bytes"
	"fmt"
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

	msg := []byte("JDWP-Handshake")
	if n, err := conn.Write(msg); n != len(msg) || err != nil {
		return nil, fmt.Errorf("read handhake error count=%d: %s\n", n, err)
	}

	var msgr [14]byte
	r, err := conn.Read(msgr[:])
	if err != nil {
		return nil, fmt.Errorf("handhake errors [read]=%d: %s\n", r, err)
	}
	if !bytes.Equal(msg, msgr[:]) {
		return nil, fmt.Errorf("invalid handhake response")
	}

	log.Printf("success handhake")
	return &Conn{conn}, nil
}
