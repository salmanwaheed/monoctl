package monoctl

import (
  "fmt"
  "io/fs"
  "os"
  "path/filepath"
  "strings"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

type Config struct {
  EnvPrefix string // <PREFIX>_KEY, e.g., MONOCTL_DB_PORT=8080
  DefaultFilePath string // auto generate default value: $HOME/<EnvPrefix>/config.yml

  filePath string // points to a string variable
}

func (c *Config) expandFilePath(p string) string {
  return os.ExpandEnv(p)
}

func (c *Config) defaultFilePath() string {
  filePath := c.expandFilePath(c.DefaultFilePath)

  if filePath == "" {
    dir, _ := os.UserConfigDir()
    return fmt.Sprintf("%s/%s/config.yml", dir, strings.ToLower(c.EnvPrefix))
  }

  return filePath
}

func (c *Config) BindConfigFlag(cmd *cobra.Command) {
  cmd.PersistentFlags().StringVar(&c.filePath, "config", c.defaultFilePath(), "config file")
}

// load config
func (c *Config) Load() {
  CheckErr(c.load())
}

// load config
func (c *Config) load() error {
  cfgFile := c.expandFilePath(c.filePath)
  cfgDir := filepath.Dir(cfgFile)

  if err := os.MkdirAll(cfgDir, 0755); err != nil {
    return fmt.Errorf("cannot create directory: %v", err)
  }

  viper.SetConfigFile(cfgFile)

  viper.SetEnvPrefix(strings.ToUpper(c.EnvPrefix))
  viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 'db.host:"value"' to '<envPrefix>_DB_HOST="<value>"'
  viper.AutomaticEnv()

  if err := viper.ReadInConfig(); err != nil {
    if e, ok := err.(viper.UnsupportedConfigError); ok {
      return fmt.Errorf("%v", e)
    } else if _, ok := err.(*fs.PathError); ok {
      if err := viper.SafeWriteConfigAs(cfgFile); err != nil {
        return fmt.Errorf("unable to create a new config file: %w", err)
      }
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

// get all config as map
func (c *Config) getData() map[string]any {
  return viper.AllSettings()
}

// print key value, error if missing
func (c *Config) Get(cmd *cobra.Command, args []string) error {
  key := args[0]
  if err := c.isSet(key); err != nil {
    return err
  }

  val := viper.Get(key)
  fmt.Printf("%v\n", val)

  return nil
}

// set key and save to file
func (c *Config) Set(cmd *cobra.Command, args []string) error {
  key := args[0]
  value := args[1]

  viper.Set(key, value)

  if err := viper.WriteConfig(); err != nil {
    return fmt.Errorf("cannot write config: %v", err)
  }

  return nil
}

// remove key from file
func (c *Config) Unset(cmd *cobra.Command, args []string) error {
  key := args[0]

  if err := c.isSet(key); err != nil {
    return err
  }

  data := c.getData()
  path := strings.Split(key, ".")
  lastKey := strings.ToLower(path[len(path)-1])
  deepestMap := c.deepSearch(data, path[0:len(path)-1])

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
func (c *Config) List(cmd *cobra.Command, args []string) error {
  return executeFormatFlag(cmd, c.getData())
}

// check if key exists anywhere in Viper (file, env, default, Set)
func (c *Config) isSet(key string) error {
  key = strings.ToLower(key)

  if !viper.IsSet(key) {
    return fmt.Errorf("invalid config path: %v", key)
  }

  return nil
}

// get nested map at path, nil if missing
func (c *Config) deepSearch(m map[string]any, path []string) map[string]any {
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
