package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "monoctl",
	Short: "A Go based CLI tool",
	SilenceErrors: true,
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$HOME/.config/monoctl/config.yml", "config file")

	// config loader
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cfgFile = os.ExpandEnv(cfgFile)
	cfgDir := filepath.Dir(cfgFile)

	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		fmt.Println("CONFIG ERROR:", err)
	}

	viper.SetConfigFile(cfgFile)

	viper.SetEnvPrefix("MONOCTL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 'db.host:"value"' to 'MONOCTL_DB_HOST="<value>"'
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Printf("fatal error config file: %v\n", err)
		fmt.Println("config missing, creating:", cfgFile)
		viper.SafeWriteConfigAs(cfgFile)
	}

	fmt.Printf("Using config: %v\n", viper.ConfigFileUsed())
}
