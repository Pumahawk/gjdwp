package jdwp

import (
	"bytes"
	"fmt"
	"io"
)

type VersionCommandReply struct {
	Description string
	JdwpMajor   int32
	JdwpMinor   int32
	VmVersion   string
	VmName      string
}

func ParseVersionCommandReply(r io.Reader) (*VersionCommandReply, error) {
	errorf := "jdwp parse version command %q: %w"

	description, err := parseString(r)
	if err != nil {
		return nil, fmt.Errorf(errorf, description, err)
	}
	jdwpMajor, err := parseInt(r)
	if err != nil {
		return nil, fmt.Errorf(errorf, jdwpMajor, err)
	}
	jdwpMinor, err := parseInt(r)
	if err != nil {
		return nil, fmt.Errorf(errorf, jdwpMinor, err)
	}
	vmVersion, err := parseString(r)
	if err != nil {
		return nil, fmt.Errorf(errorf, vmVersion, err)
	}
	vmName, err := parseString(r)
	if err != nil {
		return nil, fmt.Errorf(errorf, vmName, err)
	}

	return &VersionCommandReply{
		Description: description,
		JdwpMajor:   jdwpMajor,
		JdwpMinor:   jdwpMinor,
		VmVersion:   vmVersion,
		VmName:      vmName,
	}, nil
}

func (v *VersionCommandReply) Bytes() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write(encodeString(v.Description))
	buf.Write(encodeInt(v.JdwpMajor))
	buf.Write(encodeInt(v.JdwpMinor))
	buf.Write(encodeString(v.VmVersion))
	buf.Write(encodeString(v.VmName))
	return buf.Bytes(), nil
}
