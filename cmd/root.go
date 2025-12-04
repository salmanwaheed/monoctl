package cmd

import (
	"fmt"

	"github.com/salmanwaheed/monoctl"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
  Use: "monoctl",
  Short: "A Go based CLI tool",
  SilenceErrors: true,
  SilenceUsage: true,
  CompletionOptions: cobra.CompletionOptions{
    DisableDefaultCmd: true,
  },
}

func Execute() {
  monoctl.CheckErr(rootCmd.Execute())
}

func init() {
  // global flags
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$HOME/.config/monoctl/config.yml", "config file")

  // runs once at initialization.
  // Ideal for: loading config, reading environment variables, initializing libraries, etc.
  cobra.OnInitialize(func() {
    if err := monoctl.Load(cfgFile, rootCmd.Use); err != nil {
      fmt.Printf("error: %v\n", err)
    }
  })
}
