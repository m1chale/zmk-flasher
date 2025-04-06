package main

import (
	"os"
	"runtime"
	"zmk-flasher/cmd"
	"zmk-flasher/platform"
)

func main() {

	switch runtime.GOOS {
	case "darwin":
		platform.Operations = platform.DarwinPlatformOperations{}
	case "linux":
		platform.Operations = platform.LinuxPlatformOperations{}
	default:
		println("OS not supported yet")
		os.Exit(1)
	}

	err := cmd.Execute()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
