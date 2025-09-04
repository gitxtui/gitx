package tui

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

type testModel struct {
	Model
}

func TestModel_InitialPanels(t *testing.T) {
	m := initialModel()

	if len(m.panels) != int(totalPanels) {
		t.Fatalf("expected %d panels, but got %d", totalPanels, len(m.panels))
	}

	for i, p := range m.panels {
		if p.content != "Loading..." {
			t.Errorf("panel %s content field was not initialized correctly", Panel(i).ID())
		}
	}
}

func TestModel_DynamicLayout(t *testing.T) {
	tm := newTestModel()
	expandedHeight := int(float64(tm.height-1) * 0.3)

	testCases := []struct {
		name             string
		focusOn          Panel
		panelToCheck     Panel
		expectedHeight   int
		shouldBeExpanded bool
	}{
		{"SecondaryPanel is collapsed by default", MainPanel, SecondaryPanel, 3, false},
		{"SecondaryPanel expands on focus", SecondaryPanel, SecondaryPanel, expandedHeight, true},
		{"StashPanel is collapsed by default", MainPanel, StashPanel, 3, false},
		{"StashPanel expands on focus", StashPanel, StashPanel, expandedHeight, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tm.focusedPanel = tc.focusOn
			tm.Model = tm.recalculateLayout()
			actualHeight := tm.panelHeights[tc.panelToCheck]

			if actualHeight != tc.expectedHeight {
				t.Errorf("panel height is incorrect: got %d, want %d", actualHeight, tc.expectedHeight)
			}
		})
	}
}

func TestModel_ScrollToTopOnFocus(t *testing.T) {
	tm := newTestModel()
	tm.panels[StashPanel].viewport.SetContent(strings.Repeat("line\n", 30))
	tm.panels[StashPanel].viewport.YOffset = 10 // Manually scroll down

	// Simulate changing focus to the StashPanel
	tm.focusedPanel = MainPanel
	updatedModel, _ := tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("5")})
	tm.Model = updatedModel.(Model)

	if tm.focusedPanel != StashPanel {
		t.Fatal("Focus did not change to StashPanel as expected")
	}
	if tm.panels[StashPanel].viewport.YOffset != 0 {
		t.Errorf("StashPanel did not scroll to top on focus: YOffset is %d, want 0", tm.panels[StashPanel].viewport.YOffset)
	}
}

func TestModel_ConditionalScrollbar(t *testing.T) {
	zone.NewGlobal()
	defer zone.Close()

	tm := newTestModel()
	tm.panels[StashPanel].viewport.SetContent(strings.Repeat("line\n", 30))
	tm.panels[CommitsPanel].viewport.SetContent(strings.Repeat("line\n", 30))

	t.Run("Scrollbar is hidden when StashPanel is not focused", func(t *testing.T) {
		tm.focusedPanel = MainPanel
		rendered := tm.renderPanel("Stash", 30, tm.panelHeights[StashPanel], StashPanel)
		if strings.Contains(rendered, scrollThumb) {
			t.Error("Scrollbar thumb should be hidden but was found")
		}
	})

	t.Run("Scrollbar is visible when StashPanel is focused", func(t *testing.T) {
		tm.focusedPanel = StashPanel
		rendered := tm.renderPanel("Stash", 30, tm.panelHeights[StashPanel], StashPanel)
		if !strings.Contains(rendered, scrollThumb) {
			t.Error("Scrollbar thumb should be visible but was not found")
		}
	})

	t.Run("Normal panel scrollbar is always visible if scrollable", func(t *testing.T) {
		tm.focusedPanel = MainPanel // Focus is NOT on CommitsPanel
		rendered := tm.renderPanel("Commits", 30, tm.panelHeights[CommitsPanel], CommitsPanel)
		if !strings.Contains(rendered, scrollThumb) {
			t.Error("Scrollbar thumb should be visible but was not found")
		}
	})
}

func TestModelPanelCycle(t *testing.T) {
	tm := newTestModel()
	t.Run("shift focus to next panel", func(t *testing.T) {
		tm.focusedPanel = MainPanel
		tm.nextPanel()
		assertPanel(t, tm.focusedPanel, StatusPanel)
	})
	t.Run("shift focus to previous panel", func(t *testing.T) {
		tm.focusedPanel = StatusPanel
		tm.prevPanel()
		assertPanel(t, tm.focusedPanel, MainPanel)
	})
}

func TestModel_KeyFocus(t *testing.T) {
	testCases := []struct {
		name          string
		initialPanel  Panel
		key           string
		expectedPanel Panel
	}{
		{"Focus Next with Tab", StatusPanel, "tab", FilesPanel},
		{"Focus Previous with Shift+Tab", FilesPanel, "shift+tab", StatusPanel},
		{"Direct Focus with '0'", StashPanel, "0", MainPanel},
		{"Direct Focus with '1'", MainPanel, "1", StatusPanel},
		{"Direct Focus with '2'", MainPanel, "2", FilesPanel},
		{"Direct Focus with '3'", MainPanel, "3", BranchesPanel},
		{"Direct Focus with '4'", MainPanel, "4", CommitsPanel},
		{"Direct Focus with '5'", MainPanel, "5", StashPanel},
		{"Direct Focus with '6'", MainPanel, "6", SecondaryPanel},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := initialModel()
			m.focusedPanel = tc.initialPanel
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tc.key)}
			if tc.key == "tab" {
				keyMsg.Type = tea.KeyTab
			}
			if tc.key == "shift+tab" {
				keyMsg.Type = tea.KeyShiftTab
			}

			updatedModel, _ := m.Update(keyMsg)
			assertPanel(t, updatedModel.(Model).focusedPanel, tc.expectedPanel)
		})
	}
}

func TestModel_contextualHelp(t *testing.T) {
	m := initialModel()
	keys = DefaultKeyMap()
	t.Run("Files Panel Help", func(t *testing.T) {
		m.focusedPanel = FilesPanel
		gotKeys := m.panelShortHelp()
		assertKeyBindingsEqual(t, gotKeys, keys.FilesPanelHelp())
	})
}

func TestModel_HelpToggle(t *testing.T) {
	m := initialModel()
	t.Run("toggles help on", func(t *testing.T) {
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
		if !updatedModel.(Model).showHelp {
			t.Error("showHelp should be true after pressing '?'")
		}
	})
	t.Run("toggles help off", func(t *testing.T) {
		m.showHelp = true
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
		if updatedModel.(Model).showHelp {
			t.Error("showHelp should be false after pressing '?'")
		}
	})
}

func TestModel_MouseFocus(t *testing.T) {
	zone.NewGlobal()
	defer zone.Close()

	testCases := []struct {
		name        string
		targetPanel Panel
	}{
		{"clicking on FilesPanel changes focus", FilesPanel},
		{"clicking on SecondaryPanel changes focus", SecondaryPanel},
	}

	for _, tc := range testCases {
		t.Skip("WILL FIX")
		t.Run(tc.name, func(t *testing.T) {
			tm := newTestModel()
			tm.focusedPanel = MainPanel

			zone.Scan(tm.View())
			time.Sleep(20 * time.Millisecond)

			panelZone := zone.Get(tc.targetPanel.ID())
			if panelZone.IsZero() {
				t.Fatalf("Could not find zone for %s. Is zone.Mark() used in the View?", tc.targetPanel.ID())
			}

			msg := tea.MouseMsg{
				X:      panelZone.StartX,
				Y:      panelZone.StartY,
				Button: tea.MouseButtonLeft,
				Action: tea.MouseActionRelease,
			}
			updatedModel, _ := tm.Update(msg)

			assertPanel(t, updatedModel.(Model).focusedPanel, tc.targetPanel)
		})
	}
}

func TestModel_Update_FileWatcher(t *testing.T) {
	m := initialModel()
	// Use the blank identifier _ to ignore the returned model
	_, cmd := m.Update(fileWatcherMsg{})

	if cmd == nil {
		t.Fatal("expected a command to be returned")
	}

	cmds := cmd().(tea.BatchMsg)
	// Cast totalPanels to an int for comparison
	if len(cmds) != int(totalPanels) {
		t.Errorf("expected %d commands, got %d", totalPanels, len(cmds))
	}
}

// newTestModel creates a new model with default dimensions and a calculated layout.
func newTestModel() testModel {
	m := initialModel()
	m.width = 100
	m.height = 31
	m = m.recalculateLayout()
	return testModel{m}
}

// assertPanel is a helper to compare focused panels.
func assertPanel(t *testing.T, got, want Panel) {
	t.Helper()
	if got != want {
		t.Errorf("got focused panel %v, want %v", got, want)
	}
}

// assertKeyBindingsEqual is a helper to compare two slices of key.Binding.
func assertKeyBindingsEqual(t *testing.T, got, want []key.Binding) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n\tgot \t%v\n\twant \t%v", got, want)
	}
}
