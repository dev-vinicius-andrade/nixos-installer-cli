package types

import (
	"github.com/dev-vinicius-andrade/nioscli/types/enums"
	"github.com/dev-vinicius-andrade/nioscli/types/interfaces"
)

type CommandDefinition struct {
	interfaces.CommandDefinition
}

type CobraHelpOptions struct {
	Title string
}
type NixDiskDevices struct {
	Devices []NixDiskDevice
}
type NixDiskDevice struct {
	Name       *string
	Device     *string
	Type       *enums.NixDiskDeviceType
	Content    *NixDiskDeviceContent
	Partitions *[]NixDiskDeviceContentPartition
}
type NixDiskDeviceContent struct {
	Type *enums.NixDiskDeviceContentType
}
type NixDiskDeviceContentPartition struct {
	Type    *enums.NixDiskDeviceContentPartitionType
	Size    *string
	Content *NixDiskDeviceContentPartitionContent
}
type NixDiskDeviceContentPartitionContent struct {
	Type       *enums.NixDiskDeviceContentPartitionContentType
	Format     *enums.NixDiskDeviceContentPartitionContentFormat
	MountPoint *string
}
type DiskDeviceInformation struct {
	Device string
	Size   string
}
