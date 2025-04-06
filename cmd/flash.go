package cmd

import (
	"os"
	"github.com/new-er/zmk-flasher/views"

	tea "github.com/charmbracelet/bubbletea"
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
	flashCmd.Flags().StringVarP(&leftBootloaderFile, "left", "l", "", "The bootloader file for the left controller (mutually exclusive with --left-and-right, must be used with --right)")
	flashCmd.Flags().StringVarP(&rightBootloaderFile, "right", "r", "", "The bootloader file for the right controller (mutually exclusive with --left-and-right, must be used with --left)")
	flashCmd.Flags().StringVarP(&leftAndRightBootloaderFile, "left-and-right", "a", "", "The bootloader file for both controllers (mutually exclusive with --left and --right)")
	flashCmd.MarkFlagsRequiredTogether("left", "right")
	flashCmd.MarkFlagsMutuallyExclusive("left", "left-and-right")
	flashCmd.MarkFlagsMutuallyExclusive("right", "left-and-right")
	flashCmd.MarkFlagsOneRequired("left", "right", "left-and-right")

	flashCmd.Flags().StringVarP(&leftControllerMountPoint, "left-mount", "m", "", "The mount point for the left controller. If not provided, the program will start an interactive mount attempt")
	flashCmd.Flags().StringVarP(&rightControllerMountPoint, "right-mount", "n", "", "The mount point for the right controller. If not provided, the program will start an interactive mount attempt")

	flashCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Do not copy the bootloader files to the controllers")
}

func run() {
	if leftAndRightBootloaderFile != "" {
		leftBootloaderFile = leftAndRightBootloaderFile
		rightBootloaderFile = leftAndRightBootloaderFile
	}
	if _, err := os.Stat(leftBootloaderFile); os.IsNotExist(err) {
		println("Left bootloader file does not exist")
		os.Exit(1)
	}
	if _, err := os.Stat(rightBootloaderFile); os.IsNotExist(err) {
		println("Right bootloader file does not exist")
		os.Exit(1)
	}
	_,err := tea.NewProgram(views.NewFlashView(
		leftBootloaderFile,
		rightBootloaderFile,
		leftControllerMountPoint,
		rightControllerMountPoint,
		dryRun,
	)).Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
