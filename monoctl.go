package monoctl

import "fmt"

type BuildInfo struct {
  Version string
  Commit string
  Repository string
  Maintainer string
}

func CheckErr(e error) {
  if err := e; err != nil {
    fmt.Printf("error: %v\n", err)
  }
}
