package platform

import (
	"fmt"
	"os/exec"
	"strings"

	"howett.net/plist"
)

type DarwinPlatformOperations struct{}

func (l DarwinPlatformOperations) GetBlockDevices() ([]BlockDevice, error) {
	cmd := exec.Command("diskutil", "list", "-plist")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var response DiskutilListResponse
	_, err = plist.Unmarshal(output, &response)
	if err != nil {
		return nil, err
	}

	blockDevices := []BlockDevice{}

	for _, disk := range response.AllDisksAndPartitions {
		// handle APFSVolumes
		for _, vol := range disk.APFSVolumes {
			if isValidUserVolume(vol.VolumeName, vol.MountPoint, vol.OSInternal) {
				blockDevices = append(blockDevices, BlockDevice{
					Name:        vol.DeviceIdentifier,
					Label:       vol.VolumeName,
					UUID:        vol.VolumeUUID,
					MountPoints: []string{vol.MountPoint},
				})
			}
		}

		// handle legacy Partitions
		for _, part := range disk.Partitions {
			if isValidUserVolume(part.VolumeName, part.MountPoint, false) {
				blockDevices = append(blockDevices, BlockDevice{
					Name:        part.DeviceIdentifier,
					Label:       part.VolumeName,
					UUID:        part.VolumeUUID,
					MountPoints: []string{part.MountPoint},
				})
			}
		}

		// handle "flat" volume entries
		if disk.DeviceIdentifier != "" &&
			disk.VolumeName != "" &&
			disk.MountPoint != "" &&
			len(disk.Partitions) == 0 &&
			len(disk.APFSVolumes) == 0 &&
			isValidUserVolume(disk.VolumeName, disk.MountPoint, disk.OSInternal) {

			blockDevices = append(blockDevices, BlockDevice{
				Name:        disk.DeviceIdentifier,
				Label:       disk.VolumeName,
				UUID:        "",
				MountPoints: []string{disk.MountPoint},
			})
		}
	}

	return blockDevices, nil
}

var skipNames = []string{
	"Macintosh HD", "Preboot", "Update", "Recovery", "VM", "Hardware",
	"xART", "iSCPreboot", "iOS", "Simulator", "CoreSimulator", "System",
}

func isValidUserVolume(name, mount string, osInternal bool) bool {
	if name == "" || mount == "" || osInternal {
		return false
	}
	for _, skip := range skipNames {
		if strings.Contains(name, skip) {
			return false
		}
	}
	return true
}

type DiskutilVolume struct {
	DeviceIdentifier string `plist:"DeviceIdentifier"`
	MountPoint       string `plist:"MountPoint"`
	VolumeName       string `plist:"VolumeName"`
	VolumeUUID       string `plist:"VolumeUUID"`
	OSInternal       bool   `plist:"OSInternal"`
}

type DiskutilPartition struct {
	DeviceIdentifier string `plist:"DeviceIdentifier"`
	MountPoint       string `plist:"MountPoint"`
	VolumeName       string `plist:"VolumeName"`
	VolumeUUID       string `plist:"VolumeUUID"`
}

type DiskutilDisk struct {
	DeviceIdentifier string              `plist:"DeviceIdentifier"`
	Content          string              `plist:"Content"`
	OSInternal       bool                `plist:"OSInternal"`
	Partitions       []DiskutilPartition `plist:"Partitions"`
	APFSVolumes      []DiskutilVolume    `plist:"APFSVolumes"`
	MountPoint       string              `plist:"MountPoint"`
	VolumeName       string              `plist:"VolumeName"`
	Size             uint64              `plist:"Size"`
}

type DiskutilListResponse struct {
	AllDisksAndPartitions []DiskutilDisk `plist:"AllDisksAndPartitions"`
}

func (d DarwinPlatformOperations) MountBlockDevice(device BlockDevice) (BlockDevice, error) {
	if len(device.MountPoints) > 0 {
		return device, nil
	}

	cmd := exec.Command("diskutil", "info", "-plist", "/dev/"+device.Name)
	output, err := cmd.Output()
	if err != nil {
		return device, fmt.Errorf("diskutil info failed for %s: %w", device.Name, err)
	}

	var info struct {
		MountPoint string `plist:"MountPoint"`
	}
	_, err = plist.Unmarshal(output, &info)
	if err != nil {
		return device, fmt.Errorf("failed to parse diskutil plist: %w", err)
	}

	if info.MountPoint == "" {
		return device, fmt.Errorf("no mount point found for %s", device.Name)
	}

	device.MountPoints = append(device.MountPoints, info.MountPoint)
	return device, nil
}
