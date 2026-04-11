package jdwp

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const ReplyFlag = uint8(0x80)

type baseheader struct {
	Length int32
	Id     int32
	Flags  byte
}

type commandheader struct {
	CommandSet uint8
	Command    uint8
}

type replyheader struct {
	ErrorCode int16
}

type pack struct {
	baseheader    *baseheader
	commandheader *commandheader
	replyheader   *replyheader
	data          []byte
}

func readPack(r io.Reader) (*pack, error) {
	baseheader := baseheader{}
	binary.Read(r, binary.BigEndian, &baseheader)
	switch baseheader.Flags {
	case ReplyFlag:
		return readReplyPack(r, &baseheader)
	default:
		return readCommandPack(r, &baseheader)
	}
}

func readReplyPack(r io.Reader, bh *baseheader) (*pack, error) {
	replyheader := replyheader{}
	if err := binary.Read(r, binary.BigEndian, &replyheader); err != nil {
		return nil, fmt.Errorf("jdwp readReplyPack replyheader: %w", err)
	}
	data := make([]byte, bh.Length-11)
	if err := binRead(r, &data); err != nil {
		return nil, fmt.Errorf("jdwp readReplyPack readdata: %w", err)
	}
	return &pack{bh, nil, &replyheader, data}, nil
}

func readCommandPack(r io.Reader, bh *baseheader) (*pack, error) {
	commandheader := commandheader{}
	if err := binary.Read(r, binary.BigEndian, &commandheader); err != nil {
		return nil, fmt.Errorf("jdwp readReplyPack commandheader: %w", err)
	}
	data := make([]byte, bh.Length-11)
	if err := binRead(r, &data); err != nil {
		return nil, fmt.Errorf("jdwp readReplyPack readdata: %w", err)
	}
	return &pack{bh, &commandheader, nil, data}, nil
}

func (p *pack) write(w io.Writer) error {
	if err := binWrite(w, p.baseheader); err != nil {
		return fmt.Errorf("jdwp pack write baseheader: %w", err)
	}
	switch p.baseheader.Flags {
	case ReplyFlag:
		if err := binWrite(w, p.replyheader); err != nil {
			return fmt.Errorf("jdwp pack write replyheader: %w", err)
		}
	default:
		if err := binWrite(w, p.commandheader); err != nil {
			return fmt.Errorf("jdwp pack write commandheader: %w", err)
		}
	}
	if err := binWriteBytes(w, p.data); err != nil {
		return fmt.Errorf("jdwp pack write data: %w", err)
	}
	return nil
}
