package jdwp

import (
	"fmt"
	"log"
)

func (c *Conn) SendVirtualMachineVersion() (*VirtualMachineVersionResponse, error) {
	vmv := &VirtualMachineVersion{}
	r := c.sendCommand(vmv)
	<-r.done
	if r.err != None {
		return nil, fmt.Errorf("jdwp Conn.SendVirtualMachineVersion err: %s", r.err)
	}
	vmr := &VirtualMachineVersionResponse{}
	vmr.parse(r.data)
	return vmr, nil
}

type VirtualMachineVersion struct {
}

func (v *VirtualMachineVersion) CommandSet() uint32 {
	return 1
}

func (v *VirtualMachineVersion) Command() uint32 {
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

func (v *VirtualMachineVersionResponse) parse(data []byte) {
	log.Panicf("jwdp VirtualMachineVersionResponse.parse not implemented: data %v", data)
}
