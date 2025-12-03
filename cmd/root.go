package cmd

import (
  "fmt"

  "github.com/salmanwaheed/monoctl/pkg/config"
  "github.com/spf13/cobra"
)

// todo: fix help, add version.

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
  if err := rootCmd.Execute(); err != nil {
    fmt.Printf("error: %v\n", err)
  }
}

func init() {
  // global flags
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$HOME/.config/monoctl/config.yml", "config file")

  // runs once at initialization.
  // Ideal for: loading config, reading environment variables, initializing libraries, etc.
  cobra.OnInitialize(func() {
    if err := config.Load(cfgFile, rootCmd.Use); err != nil {
      fmt.Printf("error: %v\n", err)
    }
  })
}
