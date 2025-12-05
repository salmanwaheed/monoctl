package cmd

import (
	"fmt"
	"time"

	"github.com/salmanwaheed/monoctl"
	"github.com/spf13/cobra"
)

var (
  // sheetSrc string
  sheetId string
  sheetName string
  dbUri string
  gAuth string
  dbDate int
  cat string

  googleCmd = &cobra.Command{Use: "google", Short: "Manage Google API integrations"}

  googleSheetCmd = &cobra.Command{
    Use: "sheet",
    Short: "Manage Google Sheet operations",
    Long: "Allows reading, writing, and modifying data inside a Google Sheet",
    Run: func (cmd *cobra.Command, args []string) {
      fields := []string{"name", "email", "code", "telephone", "country", "nationality", "companyName", "salary", "category", "productTitle", "dateCreated", "lastUpdated"}

      db := monoctl.NewDB(dbUri)
      if err := db.Connect(); err != nil {
        fmt.Println(err)
      }
      defer db.Disconnect()

      query := monoctl.Query{Fields: fields, Category: cat, DateCreated: time.Date(2025, 12, dbDate, 11, 0, 0, 0, time.UTC)}

      col := db.Collection("lead")
      rawDoc, err := col.Aggregate(query.Pipeline())
      if err != nil {
        fmt.Println(err)
      }

      // fmt.Println(rawDoc.ToRows(fields))

			g := monoctl.New(gAuth, "https://www.googleapis.com/auth/spreadsheets")
			sheet := g.Sheet(sheetId, sheetName)
			sheet.Rows = rawDoc.ToRows(fields)

			if err := sheet.InsertRows(g); err != nil {
				fmt.Println(err)
			}

			fmt.Println("Data inserted to Google Sheet:", sheetId)
    },
  }
)

func init() {
  // flags
  // googleSheetCmd.Flags().StringVar(&sheetSrc, "source", "", "database source")
  googleSheetCmd.Flags().StringVar(&sheetId, "sheet-id", "", "Google Sheet ID")
  googleSheetCmd.Flags().StringVar(&sheetName, "sheet-name", "", "Google Sheet tab name")
  googleSheetCmd.Flags().StringVar(&dbUri, "db-uri", "", "mongodb connection string")
  googleSheetCmd.Flags().StringVar(&gAuth, "gauth", "", "google-service-account.json file")
  googleSheetCmd.Flags().IntVar(&dbDate, "date", 0, "date")
  googleSheetCmd.Flags().StringVar(&cat, "cat", "", "category")

  googleCmd.AddCommand(googleSheetCmd)

  rootCmd.AddCommand(googleCmd)
}
