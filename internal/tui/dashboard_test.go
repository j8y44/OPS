package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func key(m Model, code rune, text string, mod tea.KeyMod) (Model, tea.Cmd) {
	next, cmd := m.Update(tea.KeyPressMsg{Code: code, Text: text, Mod: mod})
	return next.(Model), cmd
}

func TestNavigation(t *testing.T) {
	m := New()
	if m.cursor != 0 || m.selected != 0 {
		t.Fatal("expected Overview selected")
	}
	m, _ = key(m, tea.KeyDown, "", 0)
	if m.cursor != 1 || m.selected != 0 {
		t.Fatal("moving cursor should not open a section")
	}
	m, _ = key(m, tea.KeyEnter, "", 0)
	if m.selected != 1 {
		t.Fatal("Enter should open Docker")
	}
	m, _ = key(m, tea.KeyEscape, "", 0)
	if m.selected != 0 || m.cursor != 0 {
		t.Fatal("Escape should return to Overview")
	}
	m, _ = key(m, tea.KeyUp, "", 0)
	if m.cursor != 0 {
		t.Fatal("cursor moved before first section")
	}
	for i := 0; i < 10; i++ {
		m, _ = key(m, tea.KeyDown, "", 0)
	}
	if m.cursor != len(sections)-1 {
		t.Fatal("cursor moved beyond last section")
	}
	m, _ = key(m, tea.KeyTab, "", 0)
	if m.cursor != 0 {
		t.Fatal("Tab should wrap")
	}
	m, _ = key(m, tea.KeyTab, "", tea.ModShift)
	if m.cursor != len(sections)-1 {
		t.Fatal("Shift+Tab should wrap backwards")
	}
	m, _ = key(m, tea.KeyHome, "", 0)
	if m.cursor != 0 {
		t.Fatal("Home should move to first section")
	}
	m, _ = key(m, tea.KeyEnd, "", 0)
	if m.cursor != len(sections)-1 {
		t.Fatal("End should move to last section")
	}
}

func TestDirectNavigationAndHelp(t *testing.T) {
	for i := range sections {
		m, _ := key(New(), rune('1'+i), string(rune('1'+i)), 0)
		if m.selected != i || m.cursor != i {
			t.Fatalf("shortcut %d failed", i+1)
		}
		m, _ = key(m, '?', "?", 0)
		if !m.help {
			t.Fatal("help did not open")
		}
		m, _ = key(m, tea.KeyDown, "", 0)
		if m.cursor != i {
			t.Fatal("navigation changed behind help")
		}
		m, _ = key(m, tea.KeyEscape, "", 0)
		if m.help || m.selected != i {
			t.Fatal("closing help should preserve selected section")
		}
	}
}

func TestQuit(t *testing.T) {
	for _, help := range []bool{false, true} {
		for _, msg := range []tea.KeyPressMsg{{Code: 'q', Text: "q"}, {Code: 'c', Mod: tea.ModCtrl}} {
			m := New()
			m.help = help
			_, cmd := m.Update(msg)
			if cmd == nil {
				t.Fatal("quit did not return a command")
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Fatal("expected QuitMsg")
			}
		}
	}
}

func TestResponsiveFrames(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {39, 15}, {40, 16}, {60, 20}, {87, 24}, {88, 24}, {100, 30}, {120, 40}, {180, 50}} {
		for section := range sections {
			for _, help := range []bool{false, true} {
				t.Run(fmt.Sprintf("%dx%d/section%d/help%v", size[0], size[1], section, help), func(t *testing.T) {
					m := New()
					resized, _ := m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
					m = resized.(Model)
					m.selected = section
					m.cursor = section
					m.help = help
					frame := m.render()
					if size[0] == 0 || size[1] == 0 {
						if frame != "" {
							t.Fatal("expected empty frame before terminal size")
						}
						return
					}
					lines := strings.Split(frame, "\n")
					if len(lines) != size[1] {
						t.Fatalf("height %d, want %d", len(lines), size[1])
					}
					for _, line := range lines {
						if ansi.StringWidth(line) != size[0] {
							t.Fatalf("width %d, want %d", ansi.StringWidth(line), size[0])
						}
					}
					plain := ansi.Strip(frame)
					if size[0] >= 40 && size[1] >= 16 {
						if !strings.Contains(plain, "Q quit") {
							t.Fatal("quit hint missing")
						}
						if !help && !strings.Contains(plain, strings.ToUpper(sections[section])) {
							t.Fatal("section title missing")
						}
					}
				})
			}
		}
	}
}

func TestResizePreservesState(t *testing.T) {
	m := New()
	m.cursor = 3
	m.selected = 2
	m.help = true
	for _, size := range [][2]int{{120, 40}, {40, 16}, {1, 1}, {-1, -1}, {100, 30}} {
		next, _ := m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		m = next.(Model)
		if m.cursor != 3 || m.selected != 2 || !m.help {
			t.Fatal("resize lost navigation state")
		}
		if m.width < 0 || m.height < 0 {
			t.Fatal("negative dimensions")
		}
	}
}

func TestFitUnicodeAndANSI(t *testing.T) {
	frame := fit("\x1b[31m界界界\x1b[0m\nhello", 5, 3)
	for _, line := range strings.Split(frame, "\n") {
		if ansi.StringWidth(line) != 5 {
			t.Fatalf("wrong display width: %q", line)
		}
	}
}

func TestPreviewStatus(t *testing.T) {
	m := New()
	m.width = 120
	m.height = 40
	frame := ansi.Strip(m.render())
	for _, want := range []string{"UI PREVIEW", "NO LIVE CONNECTIONS", "PLANNED", "Day 3", "Day 5", "Day 6"} {
		if !strings.Contains(frame, want) {
			t.Fatalf("missing honest integration state: %s", want)
		}
	}
	if !m.View().AltScreen {
		t.Fatal("dashboard should use the alternate screen")
	}
}
