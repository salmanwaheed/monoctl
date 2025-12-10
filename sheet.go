package monoctl

import (
  "context"
  "encoding/json"
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/option"
  "google.golang.org/api/sheets/v4"
)

type rows [][]any

func (r *rows) String() string {
  b, _ := json.Marshal(r)
  return string(b)
}

func (r *rows) Set(v string) error {
  return json.Unmarshal([]byte(v), r)
}

func (r *rows) Type() string {
  return "rows"
}

type Sheet struct {
  AuthFile string

  ID string
  TabName string
  Rows rows

  ctx context.Context
  scopes []string
  srv *sheets.Service
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

  s.srv = srv
  return srv, nil
}

func (s *Sheet) InsertRows(cmd *cobra.Command, args []string) error {
  if len(s.Rows) == 0 {
    return fmt.Errorf("no rows to insert")
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
