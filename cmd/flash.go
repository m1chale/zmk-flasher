package cmd

import (
	"os"
	"strings"
	"zmk-flasher/files"
	"zmk-flasher/platform"
	"zmk-flasher/slices"

	"github.com/spf13/cobra"
)

var (
	leftBootloaderFile         string
	rightBootloaderFile        string
	leftAndRightBootloaderFile string

	leftControllerMountPoint  string
	rightControllerMountPoint string

	dryRun bool
)

var flashCmd = &cobra.Command{
	Use:   "flash",
	Short: "Flash firmware to a keyboard",

	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	flashCmd.Flags().StringVarP(&leftBootloaderFile, "left", "l", "left.u2f", "The bootloader file for the left controller")
	flashCmd.Flags().StringVarP(&rightBootloaderFile, "right", "r", "right.u2f", "The bootloader file for the right controller")
	flashCmd.Flags().StringVarP(&leftAndRightBootloaderFile, "left-and-right", "a", "left_and_right.u2f", "The bootloader file for both controllers")
	flashCmd.MarkFlagsRequiredTogether("left", "right")
	flashCmd.MarkFlagsMutuallyExclusive("left", "left-and-right")
	flashCmd.MarkFlagsMutuallyExclusive("right", "left-and-right")
	flashCmd.MarkFlagsOneRequired("left", "right", "left-and-right")

	flashCmd.Flags().StringVarP(&leftControllerMountPoint, "left-mount", "m", "[INTERACTIVE]", "The mount point for the left controller")
	flashCmd.Flags().StringVarP(&rightControllerMountPoint, "right-mount", "n", "[INTERACTIVE]", "The mount point for the right controller")

	flashCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", true, "Print the commands that would be run without actually running them")
}

func run() {
	if leftControllerMountPoint == "[INTERACTIVE]" {
		var err error
		leftControllerDevice, err := promptForControllerConnect("left")
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		if len(leftControllerDevice.MountPoints) == 0 {
			leftControllerDevice, err = platform.Os.MountBlockDevice(leftControllerDevice)
		}
		
		leftControllerMountPoint = leftControllerDevice.MountPoints[0]
	}
	if rightControllerMountPoint == "[INTERACTIVE]" {
		var err error
		rightControllerDevice, err := promptForControllerConnect("right")
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		if len(rightControllerDevice.MountPoints) == 0 {
			rightControllerDevice, err = platform.Os.MountBlockDevice(rightControllerDevice)
		}

		rightControllerMountPoint = rightControllerDevice.MountPoints[0]
	}

	println("Flashing left bootloader")
	copyFileOrDryRun(leftBootloaderFile, leftControllerMountPoint+"/bootloader.u2f")
	println("Flashing right bootloader")
	copyFileOrDryRun(rightBootloaderFile, rightControllerMountPoint+"/bootloader.u2f")
}

func promptForControllerConnect(name string) (platform.BlockDevice, error) {
	println("Please connect the ", name, " bootloader")
	var blockDevices, err = platform.Os.GetBlockDevices()
	if err != nil {
		return platform.BlockDevice{}, err
	}

	for {
		currentBlockDevices, err := platform.Os.GetBlockDevices()
		if err != nil {
			return platform.BlockDevice{}, err
		}

		addedBlockDevices := slices.GetAddedElements(blockDevices, currentBlockDevices, func(a, b platform.BlockDevice) bool {
			return strings.EqualFold(a.UUID, b.UUID)
		})
		if len(addedBlockDevices) == 0 {
			blockDevices = currentBlockDevices
			continue
		}

		device := addedBlockDevices[0]
		println("Found ", name, "bootloader ", device.Label)
		return device, nil
	}
}

func copyFileOrDryRun(src, dst string) {
	if dryRun {
		println("cp", src, dst)
	} else {
		files.CopyFile(src, dst)
	}
}
