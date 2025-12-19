package monoctl

import "github.com/spf13/cobra"

var (
  version    = "dev"
  gitCommit  = "none"
  goVersion  = "unknown" // go1.25.5
  kernelName = "unknown" // Linux - Fedora42
  arch       = "unknown" // x86_64
)

type buildInfo struct {
  Version     string `yaml:"Version"`
  GitCommit   string `yaml:"GitCommit"`
  GoVersion   string `yaml:"GoVersion"`   // go version | awk '{print $3}'
  OS          string `yaml:"OS"`          // uname --kernel-name
  Arch        string `yaml:"Arch"`        // uname --machine
  Repository  string `yaml:"Repository"`
  Maintainer  string `yaml:"Maintainer"`
}

func NewBuildInfo() *buildInfo {
  return &buildInfo{
    Version: version,
    GitCommit: gitCommit,
    GoVersion: goVersion,
    OS: kernelName,
    Arch: arch,
    Repository: "github.com/salmanwaheed/monoctl",
    Maintainer: "Salman Waheed",
  }
}

func (b *buildInfo) Show(cmd *cobra.Command, args []string) error {
  return executeFormatFlag(cmd, *b)
}
