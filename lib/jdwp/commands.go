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
	vmr := &AllThreadsResponse{}
	err := c.sendCommandParseResponse("Conn.SendAllThreads", vmv, vmr)
	if err != nil {
		return nil, err
	}
	return vmr, nil
}
