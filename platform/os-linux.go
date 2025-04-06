package platform

import (
	"encoding/json"
	"errors"
	"fmt"
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
		if lsblkBlockDevice.Label == "" {
			continue
		}
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

func (l LinuxOsOperations) MountBlockDevice(device BlockDevice) (BlockDevice, error) {
	cmd := exec.Command("udisksctl", "mount", "-b", fmt.Sprintf("/dev/%s", device.Name))
	output, err := cmd.Output()
	if err != nil {
		if len(device.MountPoints) > 0 {
			return device, nil
		}
		return device, errors.Join(errors.New("could not mount: label: "+device.Label+" name: "+device.Name+" uuid:"+device.UUID), errors.New(string(output)), err)
	}

	stringOutput := string(output)
	splitOutput := strings.Split(stringOutput, "at")
	lastSplit := splitOutput[len(splitOutput)-1]
	mountPoint := strings.TrimSpace(lastSplit)
	mountPoint = strings.ReplaceAll(mountPoint, "\n", "")
	device.MountPoints = append(device.MountPoints, mountPoint)
	return device, nil
}
