package enums

type NixDiskDeviceContentPartitionType string
type NixDiskDeviceContentPartitionContentType string
type NixDiskDeviceContentPartitionContentFormat string
type NixDiskDeviceType string
type NixDiskDeviceContentType string

const (
	disk NixDiskDeviceType = "disk"
)

var NixDiskDeviceTypes = struct {
	Disk NixDiskDeviceType
}{
	Disk: disk,
}

const (
	Gpt NixDiskDeviceContentType = "gpt"
	Mbr NixDiskDeviceContentType = "mbr"
)
const (
	Filesystem NixDiskDeviceContentPartitionContentType = "filesystem"
	Swap       NixDiskDeviceContentPartitionContentType = "swap"
)

const (
	Ext4 NixDiskDeviceContentPartitionContentFormat = "ext4"
	Xfs  NixDiskDeviceContentPartitionContentFormat = "xfs"
	Vfat NixDiskDeviceContentPartitionContentFormat = "vfat"
)
const (
	Efi  NixDiskDeviceContentPartitionType = "EF02"
	Boot NixDiskDeviceContentPartitionType = "boot"
)
