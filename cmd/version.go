package cmd

import (
  "fmt"

  "github.com/salmanwaheed/monoctl/pkg/formatter"
  "github.com/spf13/cobra"
)

type VersionInfo struct {
  Version string
  Commit string
  Repository string
  Maintainer string
}

var verCmd = &cobra.Command{
  Use: "version",
  Short: fmt.Sprintf("version for %s", rootCmd.Root().Use),
  Args: cobra.ExactArgs(0),
  RunE: func(cmd *cobra.Command, args []string) error {
    data := VersionInfo{
      Version: "dev",
      Commit: "none",
      Repository: "github.com/salmanwaheed/monoctl",
      Maintainer: "Salman Waheed",
    }

    return formatter.Execute(cmd, data)
  },
}

func init() {
  // list flags
  formatter.BindFlag(verCmd)

  rootCmd.AddCommand(verCmd)
}
