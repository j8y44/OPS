package tui

import (
	tea "charm.land/bubbletea/v2"
)

var sections = []string{"Overview", "Docker", "Pipelines", "Kubernetes", "Security", "Settings"}

type Model struct {
	cursor   int
	selected int
	width    int
	height   int
	help     bool
}

func New() Model              { return Model{} }
func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = max(0, msg.Width), max(0, msg.Height)
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		}
		if key == "?" {
			m.help = !m.help
			return m, nil
		}
		if m.help {
			if key == "esc" {
				m.help = false
			}
			return m, nil
		}
		switch key {
		case "up", "k", "left", "h":
			m.cursor = max(0, m.cursor-1)
		case "down", "j", "right", "l":
			m.cursor = min(len(sections)-1, m.cursor+1)
		case "tab":
			m.cursor = (m.cursor + 1) % len(sections)
		case "shift+tab":
			m.cursor = (m.cursor + len(sections) - 1) % len(sections)
		case "home":
			m.cursor = 0
		case "end":
			m.cursor = len(sections) - 1
		case "enter":
			m.selected = m.cursor
		case "esc":
			m.selected, m.cursor = 0, 0
		case "1", "2", "3", "4", "5", "6":
			m.cursor = int(key[0] - '1')
			m.selected = m.cursor
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}
