package monoctl

import (
	"context"
	"fmt"
	"os"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Google struct {
	AuthFile string
	Scopes   []string

	sheetsService *sheets.Service
	mu            sync.Mutex
}

func New(authFile string, scopes ...string) *Google {
	return &Google{
		AuthFile: authFile,
		Scopes:   scopes,
	}
}

func (g *Google) clientOption() (option.ClientOption, error) {
	key, err := os.ReadFile(g.AuthFile)
	if err != nil {
		return nil, fmt.Errorf("googleapi: unable to read auth file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(key, g.Scopes...)
	if err != nil {
		return nil, fmt.Errorf("googleapi: invalid json credentials: %w", err)
	}

	return option.WithHTTPClient(config.Client(context.Background())), nil
}

// Lazy service initialization (shared session)
func (g *Google) Sheets() (*sheets.Service, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.sheetsService != nil {
		return g.sheetsService, nil
	}

	opt, err := g.clientOption()
	if err != nil {
		return nil, err
	}

	srv, err := sheets.NewService(context.Background(), opt)
	if err != nil {
		return nil, err
	}

	g.sheetsService = srv
	return srv, nil
}

type Sheet struct {
	ID   string
	Name string
	Rows [][]any
}

func (g *Google) Sheet(id, name string) *Sheet {
	return &Sheet{
		ID:   id,
		Name: name,
	}
}

func (s *Sheet) InsertRows(g *Google) error {
	if len(s.Rows) == 0 {
		return fmt.Errorf("googleapi: no rows to insert")
	}

	srv, err := g.Sheets()
	if err != nil {
		return err
	}

	vr := &sheets.ValueRange{Values: s.Rows}
	_, err = srv.Spreadsheets.Values.Append(s.ID, s.Name, vr).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()

	return err
}
