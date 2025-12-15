package cmd

import (
  "github.com/salmanwaheed/monoctl"
  "github.com/spf13/cobra"
)

var (
  sheet = &monoctl.Sheet{TabName: "Sheet1"}

  sheetCmd = &cobra.Command{
    Use: "sheet",
    Short: "Manage Google Sheet operations",
    Long: "Allows reading, writing, and modifying data inside a Google Sheet",
    RunE: sheet.InsertRows,
  }
)

func init() {
  // flags
  sheetCmd.Flags().StringVar(&sheet.AuthFile, "auth", "", "Google Service Account JSON file")
  sheetCmd.Flags().StringVar(&sheet.ID, "id", "", "Google Sheet ID")
  sheetCmd.Flags().StringVar(&sheet.TabName, "tab-name", "", "Google Sheet tab name")
  sheetCmd.Flags().Var(&sheet.Rows, "rows", "Google Sheet rows as JSON array")
  sheetCmd.Flags().StringVar(&sheet.DataSource, "data-source", "", "Get rows from database")
  sheetCmd.Flags().BoolVar(&sheet.DryRun, "dry-run", false, "Print rows only")

  // // required flags
  sheetCmd.MarkFlagRequired("auth")
  sheetCmd.MarkFlagRequired("id")
  sheetCmd.MarkFlagRequired("tab-name")
  sheetCmd.MarkFlagsOneRequired("rows", "data-source")

  rootCmd.AddCommand(sheetCmd)
}
