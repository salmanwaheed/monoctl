package monoctl

import (
  "errors"
  "fmt"
  "io/fs"
  "os"
  "path/filepath"
  "strings"

  "github.com/spf13/cobra"
  "go.yaml.in/yaml/v4"
)

const (
  DIR_NAME string = "data-sources"
  FILE_EXT string = ".yml"
)

type DataSource struct {
  Type      string    `yaml:"type"`
  Name      string    `yaml:"name"`
  Uri       string    `yaml:"uri,omitempty"`
  Table     string    `yaml:"table,omitempty"`
  Limit     int64     `yaml:"limit,omitempty"`
  Fields    []string  `yaml:"fields,omitempty"`
  Query     string    `yaml:"query,omitempty"`

  Overwrite bool      `yaml:"-"`
  DryRun    bool      `yaml:"-"`
}

var allowedTypes = map[string]struct{}{
  "mongodb": {},
  "mariadb": {},
}

func (ds *DataSource) validate() error {
  if ds.Type == "" || ds.Name == "" {
    return errors.New("type and name are required")
  }

  if _, ok := allowedTypes[ds.Type]; !ok {
    return fmt.Errorf("invalid data source type: %q", ds.Type)
  }

  return nil
}

func (ds *DataSource) getFilePath() (string, error) {
  if err := ds.validate(); err != nil {
    return "", err
  }

  dir, err := configDir(DIR_NAME, ds.Type)
  if err != nil { return "", fmt.Errorf("cannot get data-sources directory: %w", err) }

  return filepath.Join(dir, ds.Name + FILE_EXT), nil
}

func (ds *DataSource) checkArgFormat(args []string) error {
  if len(args) == 0 {
    return fmt.Errorf("missing argument: expected <type>/<name>")
  }

  parts := strings.Split(args[0], "/")
  if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
    return fmt.Errorf("invalid argument %q (expected <type>/<name>)", args[0])
  }

  ds.Type = parts[0]
  ds.Name = parts[1]
  return nil
}

func (ds *DataSource) load(args []string) (string, error) {
  if err := ds.checkArgFormat(args); err != nil { return "", err }

  path, err := ds.getFilePath()
  if err != nil { return "", err }

  b, err := os.ReadFile(path)
  if err != nil {
    if errors.Is(err, os.ErrNotExist) {
      return "", fmt.Errorf("data source '%s/%s' not found", ds.Type, ds.Name)
    }
    return "", fmt.Errorf("read data source file: %w", err)
  }

  var d DataSource
  if err := yaml.Unmarshal(b, &d); err != nil {
    return "", fmt.Errorf("cannot unmarshal data source: %w", err)
  }

  rows, err := d.RunQuery()
  if err != nil {
    return "", fmt.Errorf("execute query: %w", err)
  }

  str, err := toJson(rows)
  if err != nil {
    return "", fmt.Errorf("serialize rows to JSON: %w", err)
  }

  return str, nil
}

func (ds *DataSource) Create(cmd *cobra.Command, args []string) error {
  if err := ds.checkArgFormat(args); err != nil { return err }

  path, err := ds.getFilePath()
  if err != nil { return err }

  if !ds.Overwrite {
    if _, err := os.Stat(path); err == nil {
      return fmt.Errorf("data source '%s/%s' already exists (use --overwrite)", ds.Type, ds.Name)
    }
  }

  b, err := yaml.Marshal(ds)
  if err != nil {
    return fmt.Errorf("cannot marshal data source: %w", err)
  }

  if err := os.WriteFile(path, b, 0600); err != nil {
    return fmt.Errorf("cannot write data source file: %w", err)
  }

  return nil
}

func (ds *DataSource) Delete(cmd *cobra.Command, args []string) error {
  if err := ds.checkArgFormat(args); err != nil { return err }

  path, err := ds.getFilePath()
  if err != nil { return err }

  if err := os.Remove(path); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      return fmt.Errorf("data source '%s/%s' not found", ds.Type, ds.Name)
    }
    return fmt.Errorf("unable to delete data source %w", err)
  }

  return nil
}

func (ds *DataSource) View(cmd *cobra.Command, args []string) error {
  if err := ds.checkArgFormat(args); err != nil { return err }

  path, err := ds.getFilePath()
  if err != nil { return err }

  b, err := os.ReadFile(path)
  if err != nil {
    if errors.Is(err, os.ErrNotExist) {
      return fmt.Errorf("data source '%s/%s' not found", ds.Type, ds.Name)
    }
    return fmt.Errorf("read data source file: %w", err)
  }

  if ds.DryRun {
    var d DataSource
    if err := yaml.Unmarshal(b, &d); err != nil {
      return fmt.Errorf("cannot unmarshal data source: %w", err)
    }

    rows, err := d.RunQuery()
    if err != nil {
      return fmt.Errorf("execute query: %w", err)
    }

    str, err := toJson(rows)
    if err != nil {
      return fmt.Errorf("serialize rows to JSON: %w", err)
    }

    fmt.Println(str)
    return nil
  }

  // default: show raw YAML
  fmt.Print(string(b))
  return nil
}

func (ds *DataSource) List(cmd *cobra.Command, args []string) error {
  fmt.Printf("%-8s %-18s %-20s\n", "TYPE", "NAME", "TABLE")
  fmt.Printf("%-8s %-18s %-20s\n", "----", "----", "-----")

  dir, err := configDir(DIR_NAME)
  if err != nil { return fmt.Errorf("cannot get data-sources directory: %w", err) }

  err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
    if err != nil {
      fmt.Fprintf(os.Stderr, "warning: %v\n", err)
      return nil
    }

    if d.IsDir() || !strings.HasSuffix(d.Name(), FILE_EXT) {
      return nil
    }

    b, err := os.ReadFile(path)
    if err != nil {
      fmt.Fprintf(os.Stderr, "warning: cannot read '%s': %v\n", path, err)
      return nil
    }

    var data DataSource
    if err := yaml.Unmarshal(b, &data); err != nil {
      fmt.Fprintf(os.Stderr, "warning: cannot parse YAML %s: %v\n", path, err)
      return nil
    }

    fmt.Printf("%-8s %-18s %-20s\n", data.Type, data.Name, data.Table)

    return nil
  })

  if err != nil {
    return fmt.Errorf("failed to list data-sources: %w", err)
  }

  return nil
}

func (ds *DataSource) RunQuery() ([][]any, error) {
  runner, err := ds.runner()
  if err != nil { return nil, err }

  return runner.GetRows()
}

func (ds *DataSource) runner() (Runner, error) {
  switch ds.Type {
    case "mongodb":
      return &MongoDB{
          Uri: ds.Uri,
          Collection: ds.Table,
          Fields: ds.Fields,
          QueryJSON: ds.Query,
          Limit: ds.Limit,
        }, nil
    default:
      return nil, fmt.Errorf("unsupported data source type: %s", ds.Type)
  }
}
