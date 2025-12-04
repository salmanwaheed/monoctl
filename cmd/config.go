package cmd

import (
	"github.com/salmanwaheed/monoctl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
  cfgCmd = &cobra.Command{Use: "config", Short: "Manage config"}

  cfgSetCmd = &cobra.Command{
    Use: "set <key> <value>",
    Short: "Set a config value",
    Args: cobra.ExactArgs(2),
    RunE: func (cmd *cobra.Command, args []string) error {
      return monoctl.Set(args[0], args[1])
    },
  }

  cfgGetCmd = &cobra.Command{
    Use: "get <key>",
    Short: "Get a config value",
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      return monoctl.Get(args[0])
    },
  }

  cfgUnsetCmd = &cobra.Command{
    Use: "unset <key>",
    Short: "Unset a config value",
    Args: cobra.ExactArgs(1),
    RunE: func (cmd *cobra.Command, args []string) error {
      return monoctl.Unset(args[0])
    },
  }

  cfgListCmd = &cobra.Command{
    Use: "list",
    Short: "list a config values",
    RunE: func (cmd *cobra.Command, args []string) error {
      return monoctl.Formatter.Execute(cmd, viper.AllSettings())
    },
  }
)

func init() {
  // list flags
  monoctl.Formatter.BindFlag(cfgListCmd)

  cfgCmd.AddCommand(cfgGetCmd, cfgSetCmd, cfgUnsetCmd, cfgListCmd)

  rootCmd.AddCommand(cfgCmd)
}
