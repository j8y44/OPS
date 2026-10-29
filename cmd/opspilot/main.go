package main

import (
 "fmt"
 "os"
 "github.com/j8y44/OPS/internal/cli"
)
func main() { if err := cli.Execute(); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) } }
