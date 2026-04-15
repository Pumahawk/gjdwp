package jdwp

import (
	"fmt"
)

type VirtualMachineVersion struct {
}

func (v *VirtualMachineVersion) Data() []byte {
	return nil
}

type VirtualMachineVersionResponse struct {
	Description string
	JdwpMajor   int32
	JdwpMinor   int32
	VmVersion   string
	VmName      string
}

func (c *Conn) SendVirtualMachineVersion() (*VirtualMachineVersionResponse, error) {
	vmv := VirtualMachineVersion{}
	vmr := &VirtualMachineVersionResponse{}
	err := c.sendCommandParseResponse("Conn.SendVirtualMachineVersion", 1, 1, &vmv, vmr)
	if err != nil {
		return nil, fmt.Errorf("jdwp Conn.SendVirtualMachineVersion send command: %w", err)
	}
	return vmr, nil
}

func (v *VirtualMachineVersionResponse) parse(data []byte) error {
	return fmt.Errorf("jwdp VirtualMachineVersionResponse.parse not implemented: data %v", data)
}

type AllThreads struct {
}

func (v *AllThreads) Data() []byte {
	return nil
}

type AllThreadsResponse struct {
	Threads []uint64
}

func (c *Conn) SendAllThreads() (*AllThreadsResponse, error) {
	vmv := AllThreads{}
	vmr := &AllThreadsResponse{}
	err := c.sendCommandParseResponse("Conn.SendAllThreads", 1, 4, &vmv, vmr)
	if err != nil {
		return nil, err
	}
	return vmr, nil
}

func (v *AllThreadsResponse) parse(data []byte) error {
	return Unmashal(data, v)
}

type VMResume struct {
}

func (v *VMResume) Data() []byte {
	return nil
}

func (c *Conn) SendVMResume() (*NoResponseData, error) {
	vmv := VMResume{}
	vmr := &NoResponseData{}
	err := c.sendCommandParseResponse("Conn.SendVMResume", 1, 9, &vmv, vmr)
	if err != nil {
		return nil, err
	}
	return vmr, nil
}

type VMTopLevelThreadGroups struct {
}

func (*VMTopLevelThreadGroups) Data() []byte {
	return nil
}

func (c *Conn) SendVMTopLevelThreadGroups() (*VMTopLevelThreadGroupsResponse, error) {
	// 1, 5
	cb := VMTopLevelThreadGroups{}
	cr := &VMTopLevelThreadGroupsResponse{}
	if err := c.sendCommandParseResponse("Conn.SendVMTopLevelThreadGroups", 1, 5, &cb, cr); err != nil {
		return nil, fmt.Errorf("jdwp Conn.SendVMTopLevelThreadGroups send command: %w", err)
	}
	return cr, nil
}

type VMTopLevelThreadGroupsResponse struct {
	GroupIds []uint64
}
