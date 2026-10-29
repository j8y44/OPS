package tui

import (
 "fmt"
 "strings"
 tea "charm.land/bubbletea/v2"
)

var sections = []string{"Overview", "Docker", "Pipelines", "Kubernetes", "Security", "Settings"}

type Model struct {
 cursor int
 selected int
 width int
 height int
}

func New() Model { return Model{} }
func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
 switch msg := msg.(type) {
 case tea.WindowSizeMsg:
  m.width, m.height = msg.Width, msg.Height
 case tea.KeyPressMsg:
  switch msg.String() {
  case "ctrl+c", "q": return m, tea.Quit
  case "up", "k": if m.cursor > 0 { m.cursor-- }
  case "down", "j": if m.cursor < len(sections)-1 { m.cursor++ }
  case "enter": m.selected = m.cursor
  case "esc": m.selected = 0; m.cursor = 0
  }
 }
 return m, nil
}

func (m Model) View() tea.View {
 var b strings.Builder
 b.WriteString("OPSPILOT  /  DEVOPS CONTROL CENTER\n")
 b.WriteString("================================\n\n")
 b.WriteString("WORKSPACE\n")
 for i, section := range sections {
  prefix := "  "
  if i == m.cursor { prefix = "> " }
  fmt.Fprintf(&b, "%s%s\n", prefix, section)
 }
 b.WriteString("\n--------------------------------\n")
 fmt.Fprintf(&b, "%s\n\n", strings.ToUpper(sections[m.selected]))
 if m.selected == 0 {
  b.WriteString("Welcome to OpsPilot.\nYour local-first DevOps workspace.\n\nDocker — Planned for Day 3\nCI/CD — Planned for Day 5\n")
 } else {
  b.WriteString("This module is coming in a future sprint.\n")
 }
 b.WriteString("\n--------------------------------\n")
 b.WriteString("Up/Down Navigate | Enter Open | Esc Back | Q Quit\n")
 v := tea.NewView(b.String())
 v.AltScreen = true
 return v
}
