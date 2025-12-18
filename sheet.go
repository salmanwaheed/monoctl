package monoctl

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/option"
  "google.golang.org/api/sheets/v4"
)

type rows [][]any

func (r *rows) String() string {
  if r == nil || len(*r) == 0 {
    return "[]"
  }

  b, err := json.Marshal(*r)
  if err != nil {
    return fmt.Sprintf("failed to marshal rows: %v", err)
  }

  return string(b)
}

func (r *rows) Set(v string) error {
  if v == "" {
    return errors.New("rows input is empty")
  }

  if err := json.Unmarshal([]byte(v), r); err != nil {
    return fmt.Errorf("invalid rows JSON: %v", err)
  }

  return nil
}

func (r *rows) Type() string {
  return "rows"
}

type Sheet struct {
  AuthFile string
  DataSource string
  DryRun bool

  ID string
  TabName string
  Rows rows

  ctx context.Context
  scopes []string
  // srv *sheets.Service
}

func (s *Sheet) auth() (*sheets.Service, error) {
  s.AuthFile = os.ExpandEnv(s.AuthFile)
  s.scopes = []string{sheets.SpreadsheetsScope} // https://www.googleapis.com/auth/spreadsheets
  s.ctx = context.Background()

  jsonKey, err := os.ReadFile(s.AuthFile)
  if err != nil {
    return nil, fmt.Errorf("unable to read auth file: %v", err)
  }

  config, err := google.JWTConfigFromJSON(jsonKey, s.scopes...)
  if err != nil {
    return nil, fmt.Errorf("invalid json credentials: %v", err)
  }

  httpClient := config.Client(s.ctx)
  srv, err := sheets.NewService(s.ctx, option.WithHTTPClient(httpClient))
  if err != nil {
    return nil, fmt.Errorf("unable to retrieve client: %v", err)
  }

  // s.srv = srv
  return srv, nil
}

func (s *Sheet) validate() error {
  switch {
    case s.ID == "":
      return errors.New("spreadsheet ID is required")
    case s.TabName == "":
      return errors.New("sheet tab name is required")
    case len(s.Rows) == 0:
      return errors.New("no rows provided to insert")
  }

  return nil
}

func (s *Sheet) InsertRows(cmd *cobra.Command, args []string) error {
  if s.DataSource != "" {
    ds := &DataSource{}

    str, err := ds.load([]string{s.DataSource})
    if err != nil {
      return fmt.Errorf("load datasource %w", err)
    }

    if err := s.Rows.Set(str); err != nil {
      return fmt.Errorf("parse datasource rows: %w", err)
    }
  }

  if err := s.validate(); err != nil {
    return err
  }

  if s.DryRun {
    fmt.Println(s.Rows.String())
    return nil
  }

  srv, err := s.auth()
  if err != nil {
    return err
  }

  vr := &sheets.ValueRange{Values: s.Rows}
  _, err = srv.Spreadsheets.Values.
                                  Append(s.ID, s.TabName, vr).
                                  ValueInputOption("RAW").
                                  InsertDataOption("INSERT_ROWS").
                                  Do()

  if err != nil {
    return fmt.Errorf("unable to append data: %v", err)
  }

  return nil
}
