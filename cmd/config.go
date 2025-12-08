package cmd

import (
  "github.com/salmanwaheed/monoctl"
  "github.com/spf13/cobra"
)

var (
  cfgCmd = &cobra.Command{Use: "config", Short: "Manage config"}

  cfgGetCmd = &cobra.Command{
    Use: "get <key>",
    Short: "Get a config value",
    Args: cobra.ExactArgs(1),
    RunE: cfg.Get,
  }

  cfgSetCmd = &cobra.Command{
    Use: "set <key> <value>",
    Short: "Set a config value",
    Args: cobra.ExactArgs(2),
    RunE: cfg.Set,
  }

  cfgUnsetCmd = &cobra.Command{
    Use: "unset <key>",
    Short: "Unset a config value",
    Args: cobra.ExactArgs(1),
    RunE: cfg.Unset,
  }

  cfgListCmd = &cobra.Command{
    Use: "list",
    Short: "List a config values",
    RunE: cfg.List,
  }
)

func init() {
  // list flags
  monoctl.BindFormatFlag(cfgListCmd)

  cfgCmd.AddCommand(cfgGetCmd, cfgSetCmd, cfgUnsetCmd, cfgListCmd)
  rootCmd.AddCommand(cfgCmd)
}
