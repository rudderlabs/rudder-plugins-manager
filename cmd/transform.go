package main

import (
	"github.com/rudderlabs/rudder-transformations/cmd/commands"
	"github.com/rudderlabs/rudder-transformations/plugins"
	"github.com/spf13/cobra"
)

/**
 * This is the main entry point for the transform command.
 *  make build NAME=transform
 * ./bin/transform  json -i cmd/commands/testdata/input.json -p destinations
 */
func main() {
	rootCmd := &cobra.Command{Use: "transform"}
	rootCmd.AddCommand(commands.GetJSONCmd(plugins.ManagerInstance))
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
