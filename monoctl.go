package monoctl

import (
  "encoding/json"
  "fmt"
  "os"
  "strings"
  "text/template"

  "github.com/spf13/cobra"
  "go.yaml.in/yaml/v4"
)

type BuildInfo struct {
  Version string
  Commit string
  Repository string
  Maintainer string
}

func (b *BuildInfo) Show(cmd *cobra.Command, args []string) error {
  return executeFormatFlag(cmd, b)
}

func CheckErr(msg interface{}) {
  if msg != nil {
    fmt.Fprintln(os.Stderr, "error:", msg)
  }
}

// convert value to JSON string
func toJson(v any) (string, error) {
  b, err := json.Marshal(v)

  if err != nil {
    return "", fmt.Errorf("json encode error: %v", err)
  }

  return string(b), nil
}

// convert value to YAML string
func toYaml(v any) (string, error) {
  b, err := yaml.Marshal(v)

  if err != nil {
    return "", fmt.Errorf("yaml encode error: %v", err)
  }

  return string(b), nil
}

// bind functions to Go template
var funcMap = template.FuncMap{
  "json": toJson,
  "yaml": toYaml,
}

// render Go template with helpers
func tmpl(name string, text string, data any) error {
  t, err := template.New(name).Funcs(funcMap).Parse(text)
  if err != nil {
    return fmt.Errorf("template parse error: %v", err)
  }

  if err := t.Execute(os.Stdout, data); err != nil {
    return fmt.Errorf("template execution error: %v", err)
  }

  if !strings.Contains(text, "yaml") {
    os.Stdout.Write([]byte("\n"))
  }

  return nil
}

var valFormatFlag string

// bind flag like: monoctl <command> --format
func BindFormatFlag(cmd *cobra.Command) {
  usage := `Format output using a custom template:
'json':       Print in JSON format
'yaml':       Print in YAML format
'TEMPLATE':   Print output using the given Go template.`

  cmd.Flags().StringVar(&valFormatFlag, "format", "", usage)
}

// output as json, yaml, or template
func executeFormatFlag(cmd *cobra.Command, data any) error {
  value := strings.ToLower(strings.TrimSpace(valFormatFlag))

  switch {
    case value == "{{.}}" || !cmd.Flags().Changed("format"):
      fmt.Printf("%+v\n", data)

    case strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}"):
      return tmpl("formatterTemplate", valFormatFlag, data)

    case value == "json":
      return printFn(data, toJson, &PrintOptions{NewLine: true})

    case value == "yaml":
      return printFn(data, toYaml)

    default:
      return fmt.Errorf("invalid format value: '%s'", value)
  }

  return nil
}

type PrintOptions struct { NewLine bool }

func printFn(v any, fn func(any) (string, error), opts ...*PrintOptions) error {
  out, err := fn(v)
  if err != nil { return err }

  newline := false
  if opts != nil {
    newline = opts[0].NewLine
  }

  suffix := ""
  if newline { suffix = "\n" }

  fmt.Printf("%v%s", out, suffix)
  return nil
}
