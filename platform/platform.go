package platform

import "slices"

var Operations PlatformOperations

type PlatformOperations interface {
	GetBlockDevices() ([]BlockDevice, error)
	MountBlockDevice(device BlockDevice) (BlockDevice, error)
}

type BlockDevice struct {
	UUID        string
	Name        string
	Label       string
	MountPoints []string
}

func CompareBlockDevice(left, right BlockDevice) bool{
	return left.UUID == right.UUID && left.Name == right.Name && left.Label == right.Label && slices.Equal(left.MountPoints, right.MountPoints)
}
