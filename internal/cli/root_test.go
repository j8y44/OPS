package cli

import (
 "bytes"
 "strings"
 "testing"
)

func TestAbout(t *testing.T) {
 cmd := NewRoot()
 var out bytes.Buffer
 cmd.SetOut(&out)
 cmd.SetArgs([]string{"about"})
 if err := cmd.Execute(); err != nil { t.Fatal(err) }
 if !strings.Contains(out.String(), "OpsPilot") { t.Fatal("missing project name") }
}
