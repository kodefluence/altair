package main

import (
	"github.com/spf13/cobra"
)

func main() {
	executeCommand()
}

func executeCommand() {
	rootCmd := &cobra.Command{
		Use:   "altair",
		Short: "Light Weight and Robust API Gateway.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	_ = rootCmd.Execute()
}
