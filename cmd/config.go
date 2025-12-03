package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var format string

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Manage config",
	Args: cobra.ExactArgs(1),
}

var configSetCmd = &cobra.Command{
	Use: "set <key> <value>",
	Short: "Set a config value",
	Args: cobra.ExactArgs(2),
	Run: func (cmd *cobra.Command, args []string) {
		key, val := args[0], args[1]
		viper.Set(key, val)

		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("ERROR: cannot write config: %v\n", err)
			return
		}
	},
}

var configGetCmd = &cobra.Command{
	Use: "get <key>",
	Short: "Get a config value",
	Args: cobra.ExactArgs(1),
	Run: func (cmd *cobra.Command, args []string) {
		key := args[0]

		if viper.InConfig(key) {
			val := viper.Get(key)
			fmt.Printf("%v\n", val)
		} else {
			fmt.Printf("ERROR: cannot read config: '%v' not found\n", key)
		}
	},
}

func deepSearch(m map[string]any, path []string) map[string]any {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]any)
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]any)
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]any)
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

func Unset(key string) error {
	key = strings.ToLower(key)

	if !viper.IsSet(key) {
		return fmt.Errorf("ERROR: cannot read config: '%v' not found", key)
	}

	path := strings.Split(key, ".")
	lastKey := strings.ToLower(path[len(path)-1])
	settings := viper.AllSettings()
	deepestMap := deepSearch(settings, path[0:len(path)-1])

	delete(deepestMap, lastKey)

	v := viper.New()
	if err := v.MergeConfigMap(settings); err != nil {
		return fmt.Errorf("ERROR: cannot merge config: %w", err)
	}

	if err := v.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		return fmt.Errorf("ERROR: cannot undo config: %w", err)
	}

	return nil
}

var configUnsetCmd = &cobra.Command{
	Use: "unset <key>",
	Short: "Unset a config value",
	Args: cobra.ExactArgs(1),
	Run: func (cmd *cobra.Command, args []string) {
		if err := Unset(args[0]); err != nil {
			fmt.Println(err)
		}
	},
}

var configListCmd = &cobra.Command{
	Use: "list",
	Short: "list a config values",
	Run: func (cmd *cobra.Command, args []string) {
		settings := viper.AllSettings()

		format = strings.TrimSpace(format)

		if format == "json" {
			out, err := json.Marshal(settings)

			if err != nil {
				fmt.Printf("json encode error: %v\n", err)
			}

			fmt.Println(string(out))
			return
		}

		if strings.HasPrefix(format, "{{") && strings.HasSuffix(format, "}}") {
			tmpl := template.New("config").Funcs(template.FuncMap{
				"json": func(v ...any) string {
					if len(v) == 0 { return "Did you mean?: '{{json .}}'" }
					out, err := json.Marshal(v[0])

					if err != nil {
						return fmt.Sprintf("json encode error: %v\n", err)
					}

					return string(out)
				},
			})

			// format := strings.ReplaceAll(format, "{{", "{{json ")

			t, err := tmpl.Parse(format)
			if err != nil {
				fmt.Printf("template parse error: %v\n", err)
				return
			}

			if err := t.Execute(os.Stdout, settings); err != nil {
				fmt.Printf("template execution error: %v\n", err)
				return
			}

			fmt.Println() // add newline after template output
			return
		}

		fmt.Printf("unsupported format: %v\n", format)
	},
}

func init() {
	// list flags
	configListCmd.Flags().StringVar(&format, "format", "json", "Output format: 'json' or Go template like '{{json .db}}'")

	configCmd.AddCommand(configGetCmd, configSetCmd, configUnsetCmd, configListCmd)

	rootCmd.AddCommand(configCmd)
}
