package tui

import (
 "testing"
 tea "charm.land/bubbletea/v2"
)

func TestInitialState(t *testing.T) {
 m := New()
 if m.cursor != 0 || m.selected != 0 { t.Fatal("expected Overview selected") }
}
func TestNavigation(t *testing.T) {
 m := New()
 next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
 got := next.(Model)
 if got.cursor != 1 { t.Fatalf("cursor = %d, want 1", got.cursor) }
}
func TestWindowSize(t *testing.T) {
 m := New()
 next, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
 got := next.(Model)
 if got.width != 100 || got.height != 30 { t.Fatal("resize failed") }
}
