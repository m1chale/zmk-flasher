package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zmk-flasher",
	Short: "A tool to flash ZMK firmware to a keyboard",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()
	rootCmd.AddCommand(flashCmd)
}
