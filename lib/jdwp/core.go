package jdwp

import (
	"fmt"
	"io"
	"log"
)

type Conn struct {
	rwc         io.ReadWriteCloser
	events      eventsChannel
	sendCommand sendCommandChannel
	replies     repliesChannel
	stop        stopChannel
	done        doneChannel
}

type eventsChannel chan Event
type sendCommandChannel chan Command
type stopChannel chan any
type doneChannel chan any
type repliesChannel chan any
type Done chan any

func (c *Conn) Close() error {
	return c.rwc.Close()
}

func (c *Conn) stopRoutines() {
	select {
	case c.stop <- 1:
		log.Printf("jdwp Conne.stopRoutines: send")
	default:
		log.Printf("jdwp Conne.stopRoutines: ignored")
	}
}

func NewConn(rwc io.ReadWriteCloser) Conn {
	return Conn{
		rwc,
		make(eventsChannel),
		make(sendCommandChannel),
		make(repliesChannel),
		make(stopChannel, 10),
		make(doneChannel),
	}
}

func (c *Conn) Start() (<-chan any, error) {
	go c.readStreamRoutine()
	go c.closeRoutine()
	if err := c.handshake(); err != nil {
		return nil, fmt.Errorf("jdwc Conn.Start handshake: %s", err)
	}
	return c.done, nil
}

type CommandResponse struct {
	Done  <-chan any
	Error Error
}

type Event struct {
}

type Command interface {
	CommandSet() uint32
	Command() uint32
	Data() []byte
}

type BaseCommand struct {
	commandSet uint32
	command    uint32
}

func (b *BaseCommand) CommandSet() uint32 {
	return b.commandSet
}

func (b *BaseCommand) Command() uint32 {
	return b.command
}

type packBuffer struct {
	Length uint32
	Head   []byte
	Data   []byte
}
