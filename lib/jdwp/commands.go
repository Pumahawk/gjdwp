package jdwp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (c *Conn) SendVirtualMachineVersion() (*VirtualMachineVersionResponse, error) {
	vmv := &VirtualMachineVersion{}
	r := c.sendCommand(vmv)
	<-r.done
	if r.err != None {
		return nil, fmt.Errorf("jdwp Conn.SendVirtualMachineVersion err: %s", r.err)
	}
	vmr := &VirtualMachineVersionResponse{}
	if err := vmr.parse(r.data); err != nil {
		return nil, fmt.Errorf("jdwp Conn.SendVirtualMachineVersion parse: %w", err)
	}
	return vmr, nil
}

type VirtualMachineVersion struct {
}

func (v *VirtualMachineVersion) CommandSet() uint8 {
	return 1
}

func (v *VirtualMachineVersion) Command() uint8 {
	return 1
}

func (v *VirtualMachineVersion) Data() []byte {
	return nil
}

type VirtualMachineVersionResponse struct {
	BaseCommandResponse
	Description string
	JdwpMajor   int32
	JdwpMinor   int32
	VmVersion   string
	VmName      string
}

func (v *VirtualMachineVersionResponse) parse(data []byte) error {
	return fmt.Errorf("jwdp VirtualMachineVersionResponse.parse not implemented: data %v", data)
}

type AllThreads struct {
}

func (v *AllThreads) CommandSet() uint8 {
	return 1
}

func (v *AllThreads) Command() uint8 {
	return 4
}

func (v *AllThreads) Data() []byte {
	return nil
}

type AllThreadsResponse struct {
	BaseCommandResponse
	Threads []uint64
}

func (c *Conn) SendAllThreads() (*AllThreadsResponse, error) {
	vmv := &AllThreads{}
	r := c.sendCommand(vmv)
	<-r.done
	if r.err != None {
		return nil, fmt.Errorf("jdwp Conn.SendAllThreads err: %s", r.err)
	}
	vmr := &AllThreadsResponse{}
	if err := vmr.parse(r.data); err != nil {
		return nil, fmt.Errorf("jdwp Conn.SendAllThreads parse: %w", err)
	}
	return vmr, nil
}

func (v *AllThreadsResponse) parse(data []byte) error {
	bf := bytes.NewBuffer(data)
	var n uint32
	if err := binary.Read(bf, binary.BigEndian, &n); err != nil {
		return fmt.Errorf("error: jdwp AllThreadsResponse.parse: decode i")
	}
	ids := make([]uint64, 0, n)
	var idbuf uint64
	for i := range n {
		if err := binary.Read(bf, binary.BigEndian, &idbuf); err != nil {
			return fmt.Errorf("error: jdwp AllThreadsResponse.parse: decode threadId (%d/%d): %w", i, n, err)
		}
		ids = append(ids, idbuf)
	}
	v.Threads = ids
	return nil
}
