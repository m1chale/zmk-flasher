package stepchain

import "github.com/new-er/zmk-flasher/platform"

type MountBootloaderStep struct {
	BootloaderName       string
	previousBlockDevices []platform.BlockDevice
	addedBlockDevices    []platform.BlockDevice
}

func NewMountBootloaderStep(bootloaderName string) MountBootloaderStep {
	return MountBootloaderStep{
		BootloaderName:       bootloaderName,
		previousBlockDevices: nil,
	}
}

func (m MountBootloaderStep) GetName() string {
	return "Waiting for " + m.BootloaderName + " bootloader connection"
}

func (m MountBootloaderStep) Update() (Step, error) {
	blockDevices, err := platform.Operations.GetBlockDevices()
	if err != nil {
		return nil, err
	}

	if m.previousBlockDevices == nil || len(blockDevices) <= len(m.previousBlockDevices) {
		m.previousBlockDevices = blockDevices
		return m, nil
	}



	return m, nil
}

