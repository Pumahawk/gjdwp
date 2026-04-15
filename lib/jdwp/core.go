package jdwp

import (
	"fmt"
	"io"
	"sync"
)

type LoggerFunc = func(format string, v ...any)

type parsable interface {
	parse([]byte) error
}

type Conn struct {
	rwc        io.ReadWriteCloser
	idStore    idStore
	events     eventsChannel
	commands   commandsChannel
	replies    repliesChannel
	stop       stopChannel
	done       doneChannel
	LoggerFunc LoggerFunc
}

type eventsChannel chan Event
type commandsChannel chan sendCommandType
type stopChannel chan any
type doneChannel chan any
type repliesChannel chan any
type Done chan any

type sendCommandType struct {
	id uint32
	cm Command
}

func (c *Conn) log(format string, v ...any) {
	if c.LoggerFunc != nil {
		c.LoggerFunc(format, v...)
	}
}

func (c *Conn) Close() error {
	c.closeRoutine()
	return nil
}

func (c *Conn) stopRoutines() {
	select {
	case c.stop <- 1:
		c.log("jdwp Conne.stopRoutines: send")
	default:
		c.log("jdwp Conne.stopRoutines: ignored")
	}
}

func NewConn(rwc io.ReadWriteCloser) Conn {
	return Conn{
		rwc,
		idStore{sync.Mutex{}, 0, make(storeType)},
		make(eventsChannel),
		make(commandsChannel),
		make(repliesChannel),
		make(stopChannel, 10),
		make(doneChannel),
		nil,
	}
}

func (c *Conn) Start() (<-chan any, error) {
	if err := c.handshake(); err != nil {
		return nil, fmt.Errorf("jdwc Conn.Start handshake: %s", err)
	}
	go c.readStreamRoutine()
	go c.closeRoutine()
	go c.sendCommandRoutine()
	return c.done, nil
}

type commandResponse struct {
	done chan any
	err  Error
	data []byte
}

type Event struct {
}

type Command interface {
	CommandSet() uint8
	Command() uint8
	Data() []byte
}

type BaseCommand struct {
	commandSet uint8
	command    uint8
}

func (b *BaseCommand) CommandSet() uint8 {
	return b.commandSet
}

func (b *BaseCommand) Command() uint8 {
	return b.command
}

type packBuffer struct {
	Id   uint32
	Err  uint16
	Data []byte
}

type storeType map[uint32]*commandResponse

type idStore struct {
	m         sync.Mutex
	idCounter uint32
	store     storeType
}

func (i *idStore) Add() (uint32, *commandResponse) {
	i.m.Lock()
	defer i.m.Unlock()
	i.idCounter++
	id := i.idCounter
	cr := &commandResponse{make(Done), 0, nil}
	i.store[id] = cr
	return id, cr
}

func (i *idStore) GetAndDelete(id uint32) *commandResponse {
	i.m.Lock()
	defer i.m.Unlock()
	cr := i.store[id]
	delete(i.store, id)
	return cr
}

func (c *Conn) sendCommand(cm Command) (*commandResponse, error) {
	id, r := c.idStore.Add()
	select {
	case <-c.done:
		return nil, fmt.Errorf("jdwp Conn.sendCommand conn is done")
	case c.commands <- sendCommandType{id, cm}:
	}
	return r, nil
}

func (c *Conn) sendCommandParseResponse(msg string, cm Command, rc any) error {
	r, err := c.sendCommand(cm)
	if err != nil {
		return fmt.Errorf("jdwp %s send command: %w", msg, err)
	}
	select {
	case <-c.done:
		return fmt.Errorf("jdwp %s conn is done", msg)
	case <-r.done:
	}
	if r.err != None {
		return fmt.Errorf("jdwp %s err: %s", msg, r.err)
	}
	if err := Unmashal(r.data, rc); err != nil {
		return fmt.Errorf("jdwp %s parse: %w", msg, err)
	}
	return nil
}

type NoResponseData struct {
}

func (v *NoResponseData) parse(data []byte) error {
	return nil
}
