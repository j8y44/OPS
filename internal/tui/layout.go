package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var (
	accent  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	muted   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	active  = lipgloss.NewStyle().Bold(true).Reverse(true)
	planned = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)

// fit keeps every rendered frame within the terminal's cell dimensions, even
// while resizing. Truncation understands ANSI sequences and wide characters.
func fit(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	lines := strings.Split(content, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for i, line := range lines {
		line = ansi.Truncate(line, width, "")
		lines[i] = line + strings.Repeat(" ", max(0, width-ansi.StringWidth(line)))
	}
	for len(lines) < height {
		lines = append(lines, strings.Repeat(" ", width))
	}
	return strings.Join(lines, "\n")
}

func panel(title, content string, width, height int) string {
	inner := fit(accent.Render(title)+"\n\n"+content, width-4, height-2)
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("60")).Padding(0, 1).Render(inner)
}

func (m Model) render() string {
	w, h := m.width, m.height
	if w <= 0 || h <= 0 {
		return ""
	}
	if w < 40 || h < 16 {
		return fit("OPSPILOT\nResize to at least 40 x 16.\nQ / Ctrl+C quit", w, h)
	}
	header := accent.Render("OPS") + " / OPSPILOT" + muted.Render("    DEVOPS CONTROL CENTER")
	status := muted.Render("LOCAL WORKSPACE  /  UI PREVIEW  /  NO LIVE CONNECTIONS")
	if w < 64 {
		status = muted.Render("UI PREVIEW / NO LIVE CONNECTIONS")
	}
	bodyHeight := h - 4
	var body string
	if m.help {
		body = panel("KEYBOARD HELP", helpContent(), w, bodyHeight)
	} else if w >= 88 && h >= 24 {
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			m.sidebar(24, bodyHeight), " ", m.content(w-25, bodyHeight))
	} else {
		nav := fmt.Sprintf("%d/%d > %s  |  Enter open", m.cursor+1, len(sections), sections[m.cursor])
		body = active.Render(fit(nav, w, 1)) + "\n" + m.content(w, bodyHeight-1)
	}
	footer := "Arrows/Tab move  Enter open  1-6 jump  ? help  Q quit"
	if w < 64 {
		footer = "Tab move | Enter open | ? help | Q quit"
	}
	if m.help {
		footer = "? / Esc close help | Q quit"
	}
	return fit(header+"\n"+status+"\n"+body+"\n"+muted.Render(footer)+"\n"+
		muted.Render("OpsPilot / Day 2"), w, h)
}

func (m Model) sidebar(width, height int) string {
	rows := []string{}
	for i, name := range sections {
		marker := " "
		if i == m.selected {
			marker = "*"
		}
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		row := fmt.Sprintf("%s%s %d %s", cursor, marker, i+1, name)
		if i == m.cursor {
			row = active.Render(fit(row, width-4, 1))
		}
		rows = append(rows, row, "")
	}
	rows = append(rows, muted.Render("* Open section"), muted.Render("> Cursor"))
	return panel("WORKSPACE", strings.Join(rows, "\n"), width, height)
}

func (m Model) content(width, height int) string {
	if m.selected == 0 {
		return overview(width, height)
	}
	return panel(strings.ToUpper(sections[m.selected]), moduleContent(m.selected, width-4), width, height)
}

func overview(width, height int) string {
	intro := accent.Render("WORKSPACE OVERVIEW") + "\n" + muted.Render("Your local-first DevOps workspace.")
	cards := []struct{ title, detail string }{
		{"DOCKER", "Day 3 / Containers"},
		{"PIPELINES", "Day 5 / GitHub Actions"},
		{"DIAGNOSTICS", "Day 6 / Doctor"},
	}
	var rendered []string
	if width >= 75 {
		cardWidth := (width - 2) / 3
		for _, card := range cards {
			rendered = append(rendered, panel(card.title, planned.Render("PLANNED")+"\n"+card.detail, cardWidth, 6))
		}
		intro += "\n\n" + lipgloss.JoinHorizontal(lipgloss.Top, rendered[0], " ", rendered[1], " ", rendered[2])
	} else {
		for _, card := range cards {
			// Compact status rows leave room for guidance in short terminals.
			intro += "\n\n" + card.title + " / " + planned.Render("PLANNED") + "\n" + muted.Render(card.detail)
		}
	}
	if height >= 19 {
		intro += "\n\n" + accent.Render("NEXT UP") + "\n" +
			"Connect Docker and discover local containers on Day 3.\n\n" +
			muted.Render("No services have been queried. Status values are not live.") + "\n" +
			muted.Render("Press ? for shortcuts or choose a section to explore.")
	}
	return fit(intro, width, height)
}

func moduleContent(section, width int) string {
	var body string
	switch section {
	case 1:
		body = "PLANNED / DAY 3\n\nDocker is not connected.\n\nContainer discovery arrives on Day 3. Logs, inspection and refresh follow on Day 4.\n\nNo Docker commands have been executed."
	case 2:
		body = "PLANNED / DAY 5\n\nGitHub Actions is not connected.\n\nWorkflow runs and pipeline status arrive on Day 5.\n\nNo GitHub token is needed for this UI preview."
	case 3:
		body = "FUTURE MODULE\n\nKubernetes monitoring is outside the current 8-day sprint.\n\nNo cluster connection has been attempted."
	case 4:
		body = "PLANNED / DAY 6-7\n\nDiagnostics and JSON output arrive on Day 6. Security hardening follows on Day 7.\n\nNo security checks have run."
	case 5:
		body = "WORKSPACE INFORMATION\n\nMode: local UI preview\nIntegrations: not connected\nConfiguration: no settings saved\n\nUse ? for keyboard controls. These values are read-only during Day 2."
	}
	return ansi.Wrap(body, width, "")
}

func helpContent() string {
	return "Arrows / h j k l  Move cursor\nTab / Shift+Tab   Cycle sections\nHome / End       First / last section\nEnter            Open highlighted section\n1-6              Open section directly\nEsc              Return to Overview\n?                Toggle this help\nQ / Ctrl+C       Quit\n\nWhile help is open, navigation is paused.\nEsc closes help and preserves your place."
}
