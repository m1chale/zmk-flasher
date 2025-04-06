package platform

var Os PlatformOperations

type PlatformOperations interface {
	GetBlockDevices() ([]BlockDevice, error)
	MountBlockDevice(device BlockDevice) (BlockDevice, error)
}

type BlockDevice struct {
	UUID string
	Name string
	Label string
	MountPoints []string
}
