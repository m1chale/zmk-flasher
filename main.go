package main

import (
	"os"
	"zmk-flasher/cmd"
	"zmk-flasher/platform"
)

func main() {
	platform.Os = platform.LinuxPlatformOperations{}
	err := cmd.Execute()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

