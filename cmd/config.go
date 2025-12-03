package cmd

import (
	"github.com/salmanwaheed/monoctl/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Manage config",
	Args: cobra.ExactArgs(1),
}

var configSetCmd = &cobra.Command{
	Use: "set <key> <value>",
	Short: "Set a config value",
	Args: cobra.ExactArgs(2),
	RunE: func (cmd *cobra.Command, args []string) error {
		return config.Set(args[0], args[1])
	},
}

var configGetCmd = &cobra.Command{
	Use: "get <key>",
	Short: "Get a config value",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.Get(args[0])
	},
}

var configUnsetCmd = &cobra.Command{
	Use: "unset <key>",
	Short: "Unset a config value",
	Args: cobra.ExactArgs(1),
	RunE: func (cmd *cobra.Command, args []string) error {
		return config.Unset(args[0])
	},
}

var format string
var configListCmd = &cobra.Command{
	Use: "list",
	Short: "list a config values",
	RunE: func (cmd *cobra.Command, args []string) error {
		return config.List(format)
	},
}

func init() {
	// list flags
	configListCmd.Flags().StringVar(&format, "format", "json", "Output format: 'json' or Go template like '{{json .db}}'")

	configCmd.AddCommand(configGetCmd, configSetCmd, configUnsetCmd, configListCmd)

	rootCmd.AddCommand(configCmd)
}
