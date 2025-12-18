package cmd

import (
  "fmt"

  "github.com/salmanwaheed/monoctl"
  "github.com/spf13/cobra"
)

var (
  ver = &monoctl.BuildInfo{
    Version: "dev",
    GitCommit: "none",
    Repository: "github.com/salmanwaheed/monoctl",
    Maintainer: "Salman Waheed",
    GoVersion: "go1.25.5",  // go version | awk '{print $3}'
    OS: "linux",            // uname --kernel-name
    Arch: "x86_64",         // uname --machine
  }

  verCmd = &cobra.Command{
    Use: "version",
    Short: fmt.Sprintf("version for %s", rootCmd.Root().Use),
    Args: cobra.ExactArgs(0),
    RunE: ver.Show,
  }
)

func init() {
  // list flags
  monoctl.BindFormatFlag(verCmd)

  rootCmd.AddCommand(verCmd)
}
