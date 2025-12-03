package cmd

import (
  "fmt"

  "github.com/salmanwaheed/monoctl/internal/build"
  "github.com/spf13/cobra"
)

var verCmd = &cobra.Command{
  Use: "version",
  Short: fmt.Sprintf("version for %s", rootCmd.Root().Use),
  Args: cobra.ExactArgs(0),
  Run: func(cmd *cobra.Command, args []string) {
    t := "%s: %s, commit: %s, (%s)\n"
    n := cmd.CommandPath() // cli-name + command
    v := build.Version
    c := build.Commit
    a := build.Author

    fmt.Printf(t, n, v, c, a)
  },
}

func init() {
  rootCmd.AddCommand(verCmd)
}
