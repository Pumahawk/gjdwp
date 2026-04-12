package jdwp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type VirtualMachineVersion struct {
	BaseCommand
}

func NewVirtualMachineVersion() VirtualMachineVersion {
	return VirtualMachineVersion{BaseCommand{1, 1}}
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

func (c *Conn) SendVirtualMachineVersion() (*VirtualMachineVersionResponse, error) {
	vmv := NewVirtualMachineVersion()
	vmr := &VirtualMachineVersionResponse{}
	err := c.sendCommandParseResponse("Conn.SendVirtualMachineVersion", &vmv, vmr)
	if err != nil {
		return nil, fmt.Errorf("jdwp Conn.SendVirtualMachineVersion send command: %w", err)
	}
	return vmr, nil
}

func (v *VirtualMachineVersionResponse) parse(data []byte) error {
	return fmt.Errorf("jwdp VirtualMachineVersionResponse.parse not implemented: data %v", data)
}

type AllThreads struct {
	BaseCommand
}

func NewAllThreads() AllThreads {
	return AllThreads{BaseCommand{1, 4}}
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
	vmr := &AllThreadsResponse{}
	err := c.sendCommandParseResponse("Conn.SendAllThreads", vmv, vmr)
	if err != nil {
		return nil, err
	}
	return vmr, nil
}

func (v *AllThreadsResponse) parse(data []byte) error {
	bf := bytes.NewBuffer(data)
	var n uint32
	if err := binary.Read(bf, binary.BigEndian, &n); err != nil {
		return fmt.Errorf("jdwp AllThreadsResponse.parse: decode i")
	}
	ids := make([]uint64, 0, n)
	var idbuf uint64
	for i := range n {
		if err := binary.Read(bf, binary.BigEndian, &idbuf); err != nil {
			return fmt.Errorf("jdwp AllThreadsResponse.parse: decode threadId (%d/%d): %w", i, n, err)
		}
		ids = append(ids, idbuf)
	}
	v.Threads = ids
	return nil
}

type VMResume struct {
	BaseCommand
}

func NewResume() VMResume {
	return VMResume{BaseCommand{1, 9}}
}

func (v *VMResume) Data() []byte {
	return nil
}

func (c *Conn) SendVMResume() (*NoResponseData, error) {
	vmv := &VMResume{}
	vmr := &NoResponseData{}
	err := c.sendCommandParseResponse("Conn.SendVMResume", vmv, vmr)
	if err != nil {
		return nil, err
	}
	return vmr, nil
}
