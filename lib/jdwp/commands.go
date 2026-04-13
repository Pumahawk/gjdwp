package jdwp

import (
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
	BaseCommandResponse `jdwp:"ignore"`
	Threads             []uint64
}

func (c *Conn) SendAllThreads() (*AllThreadsResponse, error) {
	vmv := NewAllThreads()
	vmr := &AllThreadsResponse{}
	err := c.sendCommandParseResponse("Conn.SendAllThreads", &vmv, vmr)
	if err != nil {
		return nil, err
	}
	return vmr, nil
}

func (v *AllThreadsResponse) parse(data []byte) error {
	return Unmashal(data, v)
}

type VMResume struct {
	BaseCommand
}

func NewVMResume() VMResume {
	return VMResume{BaseCommand{1, 9}}
}

func (v *VMResume) Data() []byte {
	return nil
}

func (c *Conn) SendVMResume() (*NoResponseData, error) {
	vmv := NewVMResume()
	vmr := &NoResponseData{}
	err := c.sendCommandParseResponse("Conn.SendVMResume", &vmv, vmr)
	if err != nil {
		return nil, err
	}
	return vmr, nil
}
