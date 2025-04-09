package backend

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/platform"
)

func MountBlockDeviceCmd(device platform.BlockDevice) tea.Cmd {
	return func() tea.Msg {
		mountedDevice, err := platform.Operations.MountBlockDevice(device)
		if err != nil {
			return err
		}

		return BlockDeviceMountedMsg{
			BlockDevice: mountedDevice,
		}
	}
}

type BlockDeviceMountedMsg struct {
	BlockDevice platform.BlockDevice
}
