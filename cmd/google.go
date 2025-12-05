package cmd

import (
	"fmt"

	"github.com/salmanwaheed/monoctl"
	"github.com/spf13/cobra"
)

var (
  sheetSrc string
  sheetId string
  sheetName string

  googleCmd = &cobra.Command{Use: "google", Short: "Manage Google API integrations"}

  googleSheetCmd = &cobra.Command{
    Use: "sheet",
    Short: "Manage Google Sheet operations",
    Long: "Allows reading, writing, and modifying data inside a Google Sheet",
    Run: func (cmd *cobra.Command, args []string) {
			g := monoctl.New("google-service-account.json", "https://www.googleapis.com/auth/spreadsheets")

			sheet := g.Sheet(sheetId, sheetName)
      sheet.ID = "xxxxx"
      sheet.Name = "Sheet1"
			sheet.Rows = [][]any{
				{"Name", "Email", "Age"},
				{"Alice", "alice@example.com", 25},
				{"Bob", "bob@example.com", 30},
			}

			if err := sheet.InsertRows(g); err != nil {
				fmt.Println(err)
			}

			fmt.Println("Data inserted to Google Sheet:", sheetId)
    },
  }
)

func init() {
  // flags
  googleSheetCmd.Flags().StringVar(&sheetSrc, "source", "", "database source")
  googleSheetCmd.Flags().StringVar(&sheetId, "sheet-id", "", "Google Sheet ID")
  googleSheetCmd.Flags().StringVar(&sheetName, "sheet-name", "", "Google Sheet tab name")

  googleCmd.AddCommand(googleSheetCmd)

  rootCmd.AddCommand(googleCmd)
}
