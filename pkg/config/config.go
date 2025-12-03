package config

import (
  "encoding/json"
  "fmt"
  "io/fs"
  "os"
  "path/filepath"
  "strings"
  "text/template"

  "github.com/spf13/viper"
)

// load config
func Load(path string, envPrefix string) error {
  cfgFile := os.ExpandEnv(path)
  cfgDir := filepath.Dir(cfgFile)

  if err := os.MkdirAll(cfgDir, 0755); err != nil {
    return fmt.Errorf("cannot create directory: %v", err)
  }

  viper.SetConfigFile(cfgFile)

  viper.SetEnvPrefix(strings.ToUpper(envPrefix))
  viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 'db.host:"value"' to '<envPrefix>_DB_HOST="<value>"'
  viper.AutomaticEnv()

  if err := viper.ReadInConfig(); err != nil {
    if e, ok := err.(viper.UnsupportedConfigError); ok {
      return fmt.Errorf("%v", e)
    } else if _, ok := err.(*fs.PathError); ok {
      viper.SafeWriteConfigAs(cfgFile)
      return fmt.Errorf("config missing, created new file: %v", cfgFile)
    } else {
      return fmt.Errorf("unable to read config: %v", err)
    }
  }

  // else {
  // 	fmt.Println("Using config:", viper.ConfigFileUsed())
  // }

  return nil
}

// print key value, error if missing
func Get(key string) error {
  if err := isSet(key); err != nil {
    return err
  }

  val := viper.Get(key)
  fmt.Printf("%v\n", val)

  return nil
}

// set key and save to file
func Set(key string, value any) error {
  viper.Set(key, value)

  if err := viper.WriteConfig(); err != nil {
    return fmt.Errorf("cannot write config: %v", err)
  }

  return nil
}

// remove key from file
func Unset(key string) error {
  if err := isSet(key); err != nil {
    return err
  }

  data := getData()
  path := strings.Split(key, ".")
  lastKey := strings.ToLower(path[len(path)-1])
  deepestMap := deepSearch(data, path[0:len(path)-1])

  delete(deepestMap, lastKey)

  v := viper.New()
  if err := v.MergeConfigMap(data); err != nil {
    return fmt.Errorf("cannot merge config: %v", err)
  }

  if err := v.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
    return fmt.Errorf("cannot undo config: %v", err)
  }

  return nil
}

// show config (json or template)
func List(f string) error {
  format := strings.TrimSpace(f)
  data := getData()

  // --format json
  if format == "json" {
    fmt.Println(toJSON(data))

  // --format '{{json .}}' or '{{.}}'
  } else if strings.HasPrefix(format, "{{") && strings.HasSuffix(format, "}}") {
    funcMaps := template.FuncMap{
      "json": func(v ...any) string {
        if len(v) == 0 {
          return "Did you mean?: '{{json .}}'"
        }

        return toJSON(v[0])
      },
    }

    tmpl := template.New("config").Funcs(funcMaps)

    // format := strings.ReplaceAll(format, "{{", "{{json ") // auto load hack :D

    t, err := tmpl.Parse(format)
    if err != nil {
      return fmt.Errorf("template parse error: %v", err)
    }

    if err := t.Execute(os.Stdout, data); err != nil {
      return fmt.Errorf("template execution error: %v", err)
    }

    fmt.Println() // add newline after template output to remove '%'
  } else {
    return fmt.Errorf("unsupported format: %v", format)
  }

  return nil
}

// get all config as map
func getData() map[string]any {
  return viper.AllSettings()
}

// convert data to JSON string
func toJSON(v any) string {
  out, err := json.Marshal(v)

  if err != nil {
    return fmt.Sprintf("json encode error: %v", err)
  }

  return string(out)
}

// check if key exists anywhere in Viper (file, env, default, Set)
func isSet(key string) error {
  key = strings.ToLower(key)

  if !viper.IsSet(key) {
    return fmt.Errorf("invalid config path: %v", key)
  }

  return nil
}

// get nested map at path, nil if missing
func deepSearch(m map[string]any, path []string) map[string]any {
  for _, k := range path {
    m2, ok := m[k]
    if !ok {
      return nil // missing parent
    }

    m3, ok := m2.(map[string]any)
    if !ok {
      return nil // parent is not a map
    }

    m = m3 // continue search from here
  }
  return m
}
