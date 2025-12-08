package cmd

import (
  "github.com/salmanwaheed/monoctl"
  "github.com/spf13/cobra"
)

var (
  cfg = &monoctl.Config{
    EnvPrefix: "monoctl",
  }

  rootCmd = &cobra.Command{
    Use: "monoctl",
    Short: "A Go based CLI tool",
    SilenceErrors: true,
    SilenceUsage: true,
    CompletionOptions: cobra.CompletionOptions{
      DisableDefaultCmd: true,
    },
  }
)

func Execute() {
  monoctl.CheckErr(rootCmd.Execute())
}

func init() {
  // global flags
  cfg.BindConfigFlag(rootCmd)

  // runs once at initialization.
  // Ideal for: loading config, reading environment variables, initializing libraries, etc.
  cobra.OnInitialize(cfg.Load)
}
