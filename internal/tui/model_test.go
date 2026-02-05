package tui

import (
	"fmt"
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
	expandedHeight := int(float64(tm.height-1) * 0.4)

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
	longContent := strings.Repeat("line\n", 30)
	longContentLines := strings.Split(longContent, "\n")

	// Set up content for StashPanel
	tm.panels[StashPanel].content = longContent
	tm.panels[StashPanel].lines = longContentLines
	tm.panels[StashPanel].viewport.SetContent(longContent)

	// Set up content for CommitsPanel
	tm.panels[CommitsPanel].content = longContent
	tm.panels[CommitsPanel].lines = longContentLines
	tm.panels[CommitsPanel].viewport.SetContent(longContent)

	t.Run("Scrollbar is hidden when StashPanel is not focused", func(t *testing.T) {
		tm.focusedPanel = MainPanel
		rendered := tm.renderPanel("Stash", 30, tm.panelHeights[StashPanel], StashPanel)
		if strings.Contains(rendered, scrollThumbChar) {
			t.Error("Scrollbar thumb should be hidden but was found")
		}
	})

	t.Run("Scrollbar is visible when StashPanel is focused", func(t *testing.T) {
		tm.focusedPanel = StashPanel
		rendered := tm.renderPanel("Stash", 30, tm.panelHeights[StashPanel], StashPanel)
		if !strings.Contains(rendered, scrollThumbChar) {
			t.Error("Scrollbar thumb should be visible but was not found")
		}
	})

	t.Run("Normal panel scrollbar is always visible if scrollable", func(t *testing.T) {
		tm.focusedPanel = MainPanel // Focus is NOT on CommitsPanel
		rendered := tm.renderPanel("Commits", 30, tm.panelHeights[CommitsPanel], CommitsPanel)
		if !strings.Contains(rendered, scrollThumbChar) {
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

func TestModel_ScrollInactivePanelWithMouse(t *testing.T) {
	t.Skip("WILL FIX")
	zone.NewGlobal()
	zone.SetEnabled(true)
	defer zone.Close()

	tm := newTestModel()
	scrollablePanel := CommitsPanel
	focusedPanel := MainPanel

	// 1. Setup: Make a panel scrollable and focus another panel.
	longContent := strings.Repeat("line\n", 30)
	tm.panels[scrollablePanel].content = longContent
	tm.panels[scrollablePanel].lines = strings.Split(longContent, "\n")
	tm.panels[scrollablePanel].viewport.SetContent(longContent)
	tm.focusedPanel = focusedPanel

	// 2. Render the view and scan it to register the zones with the zone manager.
	view := tm.View()
	zone.Scan(view)

	// 3. Get the zone for the inactive panel we want to scroll.
	panelZone := zone.Get(scrollablePanel.ID())
	if panelZone.IsZero() {
		t.Fatalf("Could not find zone for %s. Is zone.Mark() used in the View?", scrollablePanel.ID())
	}

	// 4. Create a modern mouse scroll event that happens inside the inactive panel's zone.
	scrollMsg := tea.MouseMsg{
		Button: tea.MouseButtonWheelDown,
		X:      panelZone.StartX, // Within the zone's X bounds
		Y:      panelZone.StartY, // Within the zone's Y bounds
	}

	// 5. Act: Send the scroll message to the model.
	updatedModel, _ := tm.Update(scrollMsg)
	tm.Model = updatedModel.(Model)

	// 6. Assert: The inactive panel's viewport should have scrolled.
	if tm.panels[scrollablePanel].viewport.YOffset == 0 {
		t.Errorf("expected inactive panel %s to scroll, but YOffset remained 0", scrollablePanel.ID())
	}
	if tm.focusedPanel != focusedPanel {
		t.Errorf("focus should not have changed. want %s, got %s", focusedPanel.ID(), tm.focusedPanel.ID())
	}
}

func TestModel_LineSelectionAndScrolling(t *testing.T) {
	selectablePanels := []Panel{FilesPanel, BranchesPanel, CommitsPanel, StashPanel}

	for _, panel := range selectablePanels {
		t.Run(fmt.Sprintf("for %s", panel.ID()), func(t *testing.T) {
			tm := newTestModel()
			// Make viewport small to test scrolling behavior
			tm.panels[panel].viewport.Height = 5

			longContent := strings.Repeat("line\n", 10)
			lines := strings.Split(strings.TrimSpace(longContent), "\n")
			tm.panels[panel].lines = lines
			tm.panels[panel].content = longContent
			tm.panels[panel].viewport.SetContent(longContent)
			tm.focusedPanel = panel

			// Test mouse click selection
			t.Run("mouse click moves cursor", func(t *testing.T) {
				clickMsg := lineClickedMsg{panel: panel, lineIndex: 3}
				updatedModel, _ := tm.Update(clickMsg)
				tm.Model = updatedModel.(Model)

				if tm.panels[panel].cursor != 3 {
					t.Errorf("cursor should be at index 3 after click, got %d", tm.panels[panel].cursor)
				}
			})

			// Test keyboard down and up
			t.Run("arrow keys move cursor", func(t *testing.T) {
				tm.panels[panel].cursor = 1 // reset
				downKey := tea.KeyMsg{Type: tea.KeyDown}
				updatedModel, _ := tm.Update(downKey)
				tm.Model = updatedModel.(Model)

				if tm.panels[panel].cursor != 2 {
					t.Errorf("cursor should be at index 2 after pressing down, got %d", tm.panels[panel].cursor)
				}

				upKey := tea.KeyMsg{Type: tea.KeyUp}
				updatedModel, _ = tm.Update(upKey)
				tm.Model = updatedModel.(Model)

				if tm.panels[panel].cursor != 1 {
					t.Errorf("cursor should be at index 1 after pressing up, got %d", tm.panels[panel].cursor)
				}
			})

			// Test viewport scrolling down
			t.Run("viewport scrolls down with cursor", func(t *testing.T) {
				tm.panels[panel].cursor = 4 // Last visible line (0,1,2,3,4)
				tm.panels[panel].viewport.YOffset = 0

				downKey := tea.KeyMsg{Type: tea.KeyDown}
				updatedModel, _ := tm.Update(downKey)
				tm.Model = updatedModel.(Model)

				if tm.panels[panel].cursor != 5 {
					t.Fatalf("cursor should be at index 5, got %d", tm.panels[panel].cursor)
				}
				// cursor is 5, height is 5. YOffset should be 5 - 5 + 1 = 1
				if tm.panels[panel].viewport.YOffset != 1 {
					t.Errorf("viewport should scroll down. YOffset should be 1, got %d", tm.panels[panel].viewport.YOffset)
				}
			})

			// Test viewport scrolling up
			t.Run("viewport scrolls up with cursor", func(t *testing.T) {
				tm.panels[panel].cursor = 1
				tm.panels[panel].viewport.YOffset = 1

				upKey := tea.KeyMsg{Type: tea.KeyUp}
				updatedModel, _ := tm.Update(upKey)
				tm.Model = updatedModel.(Model)

				if tm.panels[panel].cursor != 0 {
					t.Fatalf("cursor should be at index 0, got %d", tm.panels[panel].cursor)
				}
				if tm.panels[panel].viewport.YOffset != 0 {
					t.Errorf("viewport should scroll up. YOffset should be 0, got %d", tm.panels[panel].viewport.YOffset)
				}
			})
		})
	}
}

func TestModel_Update_FileWatcher(t *testing.T) {
	m := initialModel()
	_, cmd := m.Update(fileWatcherMsg{})

	if cmd == nil {
		t.Fatal("expected a command to be returned")
	}

	cmds := cmd().(tea.BatchMsg)
	expectedCmds := 5
	if len(cmds) != expectedCmds {
		t.Errorf("expected %d commands, got %d", expectedCmds, len(cmds))
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
