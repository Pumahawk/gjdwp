package jdwp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type id uint32

type replyStorage struct {
	m     sync.Mutex
	store map[id]*messageData
}

func (r *replyStorage) Get(id id) *messageData {
	r.m.Lock()
	defer r.m.Unlock()
	return r.store[id]
}

func (r *replyStorage) Add(id id) *messageData {
	r.m.Lock()
	defer r.m.Unlock()
	m := &messageData{nil, make(chan interface{})}
	r.store[id] = m
	return m
}

type messageData struct {
	data Data
	done chan any
}

type idgenerator struct {
	m     sync.Mutex
	count id
}

func (i *idgenerator) Get() id {
	i.m.Lock()
	defer i.m.Unlock()
	i.count++
	return i.count
}

type Conn struct {
	conn        net.Conn
	idgenerator idgenerator
	events      chan *Event
	commands    chan *CommandPack
}

func Dial(a string) (*Conn, error) {
	log.Printf("start jdwp connection %q", a)
	conn, err := net.Dial("tcp", a)
	if err != nil {
		return nil, fmt.Errorf("invalid connection: %s", err)
	}
	handhake(conn)
	log.Printf("success handhake")
	events := make(chan *Event)
	commands := make(chan *CommandPack)
	return &Conn{conn, idgenerator{}, events, commands}, nil
}

func commandsRoutine(w io.Writer, commands <-chan *CommandPack) {
	for c := range commands {
		if err := writePack(w, c); err != nil {
			log.Printf("jdwp commandsRoutine Conn.Write: %w", err)
		}
	}
}

func responseRoutine(r io.Reader, rs *replyStorage) {
	for {
		p, err := nextPack(r)
		if err != nil {
			log.Printf("jdwp responseRoutine: error retrieve next pack")
		}
		switch p := p.(type) {
		case *ReplyPack:
			if md := rs.Get(id(p.Id)); md != nil {
				md.data = md.data
				close(md.done)
			}
		default:
			log.Printf("jdwp responseRoutnie: unexpected pack type %T", p)
		}
	}
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

func (c *Conn) Write(p *CommandPack) error {
	if err := writePack(c.conn, p); err != nil {
		return fmt.Errorf("jdwp Conn.Write: %w", err)
	}
	return nil
}

type ResumeCommandDataType string

const ResumeCommandData = ResumeCommandDataType("ResumeCommandDataType")

func NewResumeCommand(id uint32) *CommandPack {
	return &CommandPack{Pack{id}, 1, 9, ResumeCommandData}
}

func (r ResumeCommandDataType) JDWPData() []byte {
	return nil
}

func (c *Conn) SandVirtualMachineResume() error {
	if err := writePack(c.conn, NewResumeCommand(c.newId())); err != nil {
		return fmt.Errorf("jdwp Conn.Write: %w", err)
	}
	return nil
}

func (c *Conn) SandVersionCommand() error {
	writePack(c.conn)
}

func (c *Conn) newId() uint32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.countIds++
	return c.countIds
}
