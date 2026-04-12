package jdwp

func (c *Conn) SendVirtualMachineVersion() (VirtualMachineVersionResponse, error) {
	panic("Not implemented")
}

type VirtualMachineVersionResponse struct {
	CommandResponse
	Description string
	JdwpMajor   int32
	JdwpMinor   int32
	VmVersion   string
	VmName      string
}
