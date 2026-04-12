package jdwp

import (
	"fmt"
	"io"
)

const ReplyFlag = uint8(0x80)
const CommandFlag = uint8(0x00)

type PackType any

type Pack struct {
	Id uint32
}

type ReplyPack struct {
	Pack
	Error uint16
	Data  ReplyData
}

type CommandSet uint8
type Command uint8

type CommandPack struct {
	Pack
	CommandSet CommandSet
	Command    Command
	Data       CommandData
}

type Data any
type ReplyData any
type CommandData interface {
	JDWPData() []byte
}

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
		return &ReplyPack{Pack{id}, errorCode, data}, nil
	default:
		var commandSet CommandSet
		if err := binRead(r, &commandSet); err != nil {
			return nil, fmt.Errorf("jdwp pack read command commandSet: %w", err)
		}
		var command Command
		if err := binRead(r, &command); err != nil {
			return nil, fmt.Errorf("jdwp pack read command: %w", err)
		}
		data, err := solvePackData(r, length, commandSet, command)
		if err != nil {
			return nil, fmt.Errorf("jdwp pack read command data: %w", err)
		}
		return &CommandPack{Pack{id}, CommandSet(commandSet), Command(command), data}, nil
	}
}

func solvePackData(r io.Reader, length uint32, cs CommandSet, c Command) (CommandData, error) {
	bf := make([]byte, length-11)
	if err := binRead(r, &bf); err != nil {
		return nil, fmt.Errorf("jdwp pack solvePackData read bf: %w", err)
	}
	switch cs {
	case 64:
		switch c {
		case 100:
			return readEventComposite(bf)
		default:
		}
		return nil, fmt.Errorf("jdwp solvePackData cs=%d c=%d", cs, c)
	default:
		return nil, fmt.Errorf("jdwp solvePackData cs=%d", cs)
	}
}

func writePack(w io.Writer, c *CommandPack) error {
	data := c.Data.JDWPData()
	var length uint32 = uint32(len(data)) + 11
	if err := binWrite(w, length); err != nil {
		return fmt.Errorf("jdwp write pack length: %w", err)
	}
	if err := binWrite(w, c.Id); err != nil {
		return fmt.Errorf("jdwp write pack c.Id: %w", err)
	}
	if err := binWrite(w, CommandFlag); err != nil {
		return fmt.Errorf("jdwp write pack ReplyFlag: %w", err)
	}
	if err := binWrite(w, c.CommandSet); err != nil {
		return fmt.Errorf("jdwp write pack c.CommandSet: %w", err)
	}
	if err := binWrite(w, c.Command); err != nil {
		return fmt.Errorf("jdwp write pack c.Command: %w", err)
	}
	if err := binWrite(w, data); err != nil {
		return fmt.Errorf("jdwp write pack data: %w", err)
	}
	return nil
}
