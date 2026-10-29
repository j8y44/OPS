package cli

import (
 "fmt"
 "github.com/spf13/cobra"
 tea "charm.land/bubbletea/v2"
 "github.com/j8y44/OPS/internal/tui"
)

func NewRoot() *cobra.Command {
 root := &cobra.Command{
  Use: "opspilot",
  Short: "OpsPilot — terminal DevOps control center",
  Version: "0.1.0-dev",
  Args: cobra.NoArgs,
  RunE: func(cmd *cobra.Command, args []string) error {
   _, err := tea.NewProgram(tui.New()).Run()
   return err
  },
 }
 root.AddCommand(&cobra.Command{
  Use: "about",
  Short: "Show project information",
  Args: cobra.NoArgs,
  Run: func(cmd *cobra.Command, args []string) {
   fmt.Fprintln(cmd.OutOrStdout(), "OpsPilot: local-first DevOps terminal control center")
  },
 })
 return root
}

func Execute() error { return NewRoot().Execute() }
