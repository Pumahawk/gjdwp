package jdwp

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type Conn struct {
	rwc      io.ReadWriteCloser
	idStore  idStore
	events   eventsChannel
	commands commandsChannel
	replies  repliesChannel
	stop     stopChannel
	done     doneChannel
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

func (c *Conn) Close() error {
	c.closeRoutine()
	return nil
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
		idStore{sync.Mutex{}, 0, make(storeType)},
		make(eventsChannel),
		make(commandsChannel),
		make(repliesChannel),
		make(stopChannel, 10),
		make(doneChannel),
	}
}

func (c *Conn) Start() (<-chan any, error) {
	go c.readStreamRoutine()
	go c.closeRoutine()
	go c.sendCommandRoutine()
	if err := c.handshake(); err != nil {
		return nil, fmt.Errorf("jdwc Conn.Start handshake: %s", err)
	}
	return c.done, nil
}

type commandResponse struct {
	done chan any
	err  Error
	data []byte
}

type BaseCommandResponse struct {
	Err Error
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
