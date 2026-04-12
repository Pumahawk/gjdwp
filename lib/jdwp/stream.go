package jdwp

import (
	"bytes"
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
		var id uint32
		if err := binary.Read(c.rwc, binary.BigEndian, &id); err != nil {
			log.Printf("jdwp Conn.readStreamRoutine: unable to read id %s", err)
			break
		}
		var flag uint8
		if err := binary.Read(c.rwc, binary.BigEndian, &flag); err != nil {
			log.Printf("jdwp Conn.readStreamRoutine: unable to read flag %s", err)
			break
		}
		var errcode uint16
		if err := binary.Read(c.rwc, binary.BigEndian, &errcode); err != nil {
			log.Printf("jdwp Conn.readStreamRoutine: unable to read errcode %s", err)
			break
		}
		dataLength := length - 11
		data := make([]byte, dataLength)
		if n, err := r.Read(data); n != int(dataLength) || err != nil {
			log.Printf("jdwp Conn.readStreamRoutine: unable to read data n[%d/%d]: %s", n, dataLength, err)
			break
		}
		if flag == uint8(0x80) {
			go c.processPack(packBuffer{id, errcode, data})
		} else {
			log.Printf("jwdp Conn.readStreamRoutine: TODO... commands from debugger not supported yet")
		}
	}
	log.Printf("jdwp Conn.readStreamRoutine: end")
}

func (c *Conn) closeRoutine() {
	<-c.stop
	log.Printf("jdwp Conn.closeRoutine: start close operation")
	c.rwc.Close()
	close(c.done)
	close(c.commands)
	log.Printf("jdwp Conn.closeRoutine: end close operation")
}

func (c *Conn) processPack(pack packBuffer) {
	log.Printf("jdwp Conn.processPack: print pack %T %[1]v", pack)
	if pack.Id != 0 {
		r := c.idStore.GetAndDelete(pack.Id)
		r.err = Error(pack.Err)
		r.data = pack.Data
		close(r.done)
	} else {
		log.Printf("jwdp Conn.processPack: not supported id 0")
	}
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

func (c *Conn) sendCommandRoutine() {
	defer c.closeRoutine()
loop:
	for cm := range c.commands {
		select {
		case <-c.done:
			break loop
		default:
		}
		id := cm.id
		commandSet := cm.cm.CommandSet()
		command := cm.cm.Command()
		data := cm.cm.Data()
		lenght := uint32(len(data)) + 11
		bf := bytes.NewBuffer(make([]byte, 0, lenght))
		if err := binary.Write(bf, binary.BigEndian, lenght); err != nil {
			log.Printf("jdwp Conn.sendCommandRoutine: write lenght: %s", err)
			break
		}
		if err := binary.Write(bf, binary.BigEndian, id); err != nil {
			log.Printf("jdwp Conn.sendCommandRoutine: write id: %s", err)
			break
		}
		if err := binary.Write(bf, binary.BigEndian, int8(0)); err != nil {
			log.Printf("jdwp Conn.sendCommandRoutine: write flag: %s", err)
			break
		}
		if err := binary.Write(bf, binary.BigEndian, commandSet); err != nil {
			log.Printf("jdwp Conn.sendCommandRoutine: write commandSet: %s", err)
			break
		}
		if err := binary.Write(bf, binary.BigEndian, command); err != nil {
			log.Printf("jdwp Conn.sendCommandRoutine: write command: %s", err)
			break
		}
		if _, err := bf.Write(data); err != nil {
			log.Printf("jdwp Conn.sendCommandRoutine: write data: %s", err)
			break
		}
		if n, err := c.rwc.Write(bf.Bytes()); n != int(lenght) || err != nil {
			log.Printf("jdwp Conn.sendCommandRoutine: write all buffer (%d/%d): %s", n, lenght, err)
			break
		}
	}
	log.Println("jdwp Conn.sendCommandRoutine: end")
}

func (c *Conn) sendCommand(cm Command) *commandResponse {
	id, r := c.idStore.Add()
	c.commands <- sendCommandType{id, cm}
	return r
}
