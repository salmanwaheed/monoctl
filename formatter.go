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

// Args for format flag
type flag struct {
  Name string
  Value string // default value
  Usage string // help

  Variable string // auto set as "var format string"
}

// default values
var Formatter = flag{
  Name: "format",
  Value: "",
  Usage: `Format output using a custom template:
'json':       Print in JSON format
'yaml':       Print in YAML format
'TEMPLATE':   Print output using the given Go template.`,
}

// bind flag like: monoctl <command> --format
func (f *flag) BindFlag(cmd *cobra.Command) {
  cmd.Flags().StringVar(&f.Variable, f.Name, f.Value, f.Usage)
}

// output as json, yaml, or template
func (f *flag) Execute(cmd *cobra.Command, data any) error {
  value := strings.ToLower(strings.TrimSpace(f.Variable))

  switch {
    case value == "{{.}}" || !cmd.Flags().Changed(f.Name):
      fmt.Printf("%+v\n", data)

    case strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}"):
      return tmpl("formatterTemplate", f.Variable, data)

    case value == "json":
      out, err := toJsonString(data)
      if err != nil { return err }
      fmt.Printf("%v\n", out)

    case value == "yaml":
      out, err := toYamlString(data)
      if err != nil { return err }
      fmt.Printf("%v", out)

    default:
      return fmt.Errorf("invalid format value: '%s'", value)
  }

  return nil
}

// bind functions to Go template
var funcMap = template.FuncMap{
  "json": toJsonString,
  "yaml": toYamlString,
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

// convert value to JSON string
func toJsonString(v any) (string, error) {
  b, err := json.Marshal(v)

  if err != nil {
    return "", fmt.Errorf("json encode error: %v", err)
  }

  return string(b), nil
}

// convert value to YAML string
func toYamlString(v any) (string, error) {
  b, err := yaml.Marshal(v)

  if err != nil {
    return "", fmt.Errorf("yaml encode error: %v", err)
  }

  return string(b), nil
}
