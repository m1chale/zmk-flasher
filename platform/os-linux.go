package platform

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type LinuxOsOperations struct{}

func (l LinuxOsOperations) GetBlockDevices() ([]BlockDevice, error) {
	cmd := exec.Command("lsblk", "-lf", "--json")

	byteOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	stringOutput := string(byteOutput)
	lsblkResponse := LsblkResponse{}
	err = json.Unmarshal([]byte(stringOutput), &lsblkResponse)
	if err != nil {
		return nil, err
	}

	blockDevices := []BlockDevice{}
	for _, lsblkBlockDevice := range lsblkResponse.BlockDevices {
		blockDevices = append(blockDevices, BlockDevice{
			UUID:        lsblkBlockDevice.UUID,
			Name:        lsblkBlockDevice.Name,
			Label:       lsblkBlockDevice.Label,
			MountPoints: lsblkBlockDevice.MountPoints})
	}
	return blockDevices, nil
}

type LsblkResponse struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

type LsblkBlockDevice struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	MountPoints []string `json:"mountpoints"`
}

func (l LinuxOsOperations) MountBlockDevice(device BlockDevice) (BlockDevice, error) {
	cmd := exec.Command("udisksctl", "mount", "-b", "/dev/"+device.Name)
	output, err := cmd.Output()
	if err != nil {
		return device, err 
	}

	stringOutput := string(output)
	splitOutput := strings.Split(stringOutput, "at")
	lastSplit := splitOutput[len(splitOutput)-1]
	mountPoint := strings.TrimSpace(lastSplit)
	mountPoint = strings.ReplaceAll(mountPoint, "\n", "") 
	device.MountPoints = append(device.MountPoints, mountPoint)
	return device, nil
}
