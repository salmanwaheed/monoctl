package cmd

import (
  "github.com/salmanwaheed/monoctl"
  "github.com/spf13/cobra"
)

var (
  ds = &monoctl.DataSource{}

  dsCmd = &cobra.Command{Use: "data-source", Short: "Manage reusable data sources"}

  dsCreateCmd = &cobra.Command{
    Use: "create <type>/<name>",
    Short: "Create or Update a data source",
    Args: cobra.ExactArgs(1),
    RunE: ds.Create,
  }

  dsDeleteCmd = &cobra.Command{
    Use: "delete <type>/<name>",
    Short: "Delete a data source",
    Args: cobra.ExactArgs(1),
    RunE: ds.Delete,
  }

  dsViewCmd = &cobra.Command{
    Use: "view <type>/<name>",
    Short: "View a data source",
    Args: cobra.ExactArgs(1),
    RunE: ds.View,
  }

  dsListCmd = &cobra.Command{
    Use: "list",
    Short: "List all data sources",
    Args: cobra.ExactArgs(0),
    RunE: ds.List,
  }
)

func init() {
  // create flags
  dsCreateCmd.Flags().StringVar(&ds.Uri, "uri", "", "Data source connection URI")
  dsCreateCmd.Flags().StringVar(&ds.Table, "table", "", "Collection or table name")
  dsCreateCmd.Flags().Var(&ds.Query.Filter, "filter", "Query filter (JSON)")
  dsCreateCmd.Flags().StringSliceVar(&ds.Query.Select, "select", []string{}, "Select to fetch (ordered)")
  dsCreateCmd.Flags().Var(&ds.Query.Sort, "sort", "Sort records")
  dsCreateCmd.Flags().Int64Var(&ds.Query.Limit, "limit", 0, "Maximum number of records")
  dsCreateCmd.Flags().BoolVar(&ds.Overwrite, "overwrite", false, "Overwrite Query")

  // create flags required
  dsCreateCmd.MarkFlagRequired("uri")
  dsCreateCmd.MarkFlagRequired("table")
  dsCreateCmd.MarkFlagRequired("fields")
  dsCreateCmd.MarkFlagRequired("query")

  // view flags
  dsViewCmd.Flags().BoolVar(&ds.DryRun, "dry-run", false, "Execute & preview data")

  dsCmd.AddCommand(dsCreateCmd, dsDeleteCmd, dsViewCmd, dsListCmd)
  rootCmd.AddCommand(dsCmd)
}
