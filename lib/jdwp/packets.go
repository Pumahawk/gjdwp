package jdwp

import (
	"fmt"
	"io"
)

const ReplyFlag = uint8(0x80)

type typesizemap map[uint8]uint8

type PackType any

type Pack struct {
	Id   uint32
	Data Data
}

type ReplyPack struct {
	Pack
	Error uint16
}

type CommandPack struct {
	Pack
	CommandSet uint8
	Command    uint8
	Data       Data
}

type Data []byte

func nextPack(r io.Reader) (PackType, error) {
	var length uint32
	if err := binRead(r, &length); err != nil {
		return nil, fmt.Errorf("jdwp pack read length: %w", err)
	}
	var id uint32
	if err := binRead(r, &id); err != nil {
		return nil, fmt.Errorf("jdwp pack read id: %w", err)
	}
	var flags uint8
	if err := binRead(r, &flags); err != nil {
		return nil, fmt.Errorf("jdwp pack read flags: %w", err)
	}
	switch flags {
	case ReplyFlag:
		var errorCode uint16
		if err := binRead(r, &errorCode); err != nil {
			return nil, fmt.Errorf("jdwp pack read reply errorCode: %w", err)
		}
		data := make([]byte, length-11)
		if err := binRead(r, &data); err != nil {
			return nil, fmt.Errorf("jdwp pack read reply data: %w", err)
		}
		return &ReplyPack{Pack: Pack{id, data}, Error: errorCode}, nil
	default:
		var commandSet uint8
		if err := binRead(r, &commandSet); err != nil {
			return nil, fmt.Errorf("jdwp pack read command commandSet: %w", err)
		}
		var command uint8
		if err := binRead(r, &command); err != nil {
			return nil, fmt.Errorf("jdwp pack read command: %w", err)
		}
		data := make([]byte, length-11)
		if err := binRead(r, &data); err != nil {
			return nil, fmt.Errorf("jdwp pack read command data: %w", err)
		}
		return &CommandPack{Pack{id, data}, commandSet, command, data}, nil
	}
}
