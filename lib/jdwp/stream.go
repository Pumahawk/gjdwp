package jdwp

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"slices"
)

func (c *Conn) readStreamRoutine() {
	defer c.stopRoutines()
	var r io.Reader = c.rwc
	var bufl [4]byte
	for {
		if n, err := r.Read(bufl[:]); err != nil {
			log.Printf("jwdp Conn.readStreamRoutine: unable to read lenght n=[%d]: %s", n, err)
			break
		}
		var length uint32
		if _, err := binary.Decode(bufl[:], binary.BigEndian, &length); err != nil {
			panic(err)
		}
		head := make([]byte, 7)
		if n, err := r.Read(head); n != int(7) || err != nil {
			log.Printf("jdwp Conn.readStreamRoutine: unable to read head n[%d/7]: %s", n, err)
			break
		}
		dataLength := length - 11
		data := make([]byte, dataLength)
		if n, err := r.Read(data); n != int(dataLength) || err != nil {
			log.Printf("jdwp Conn.readStreamRoutine: unable to read data n[%d/%d]: %s", n, dataLength, err)
			break
		}
		go c.processPack(packBuffer{length, head, data})
	}
	log.Printf("jdwp Conn.readStreamRoutine: end")
}

func (c *Conn) closeRoutine() {
	<-c.stop
	log.Printf("jdwp Conn.closeRoutine: start close operation")
	c.rwc.Close()
	close(c.done)
	log.Printf("jdwp Conn.closeRoutine: end close operation")
}

func (c *Conn) processPack(pack packBuffer) {
	log.Printf("jdwp Conn.processPack: print pack %T %[1]v", pack)
}

func (c *Conn) handshake() error {
	message := []byte("JDWP-Handshake")
	if n, err := c.rwc.Write(message); n != 14 || err != nil {
		return fmt.Errorf("jdwt Conn.handshake: write handshake (%d/14): %w", n, err)
	}
	var buf [14]byte
	if n, err := c.rwc.Read(buf[:]); n != 14 || err != nil {
		return fmt.Errorf("jdwt Conn.handshake: read handshake (%d/14): %w", n, err)
	}
	if !slices.Equal(buf[:], message) {
		return fmt.Errorf("jdwt Conn.handshake: handshake response")
	}
	return nil
}
