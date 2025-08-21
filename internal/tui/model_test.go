package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelPanelCycle(t *testing.T) {
	t.Run("shift focus to next panel", func(t *testing.T) {
		m := initialModel()
		m.focusedPanel = MainPanel

		m.nextPanel()
		assertPanel(t, m.focusedPanel, StatusPanel)
	})
	t.Run("shift focus to previous panel", func(t *testing.T) {
		m := initialModel()
		m.focusedPanel = StatusPanel

		m.prevPanel()
		assertPanel(t, m.focusedPanel, MainPanel)
	})
	t.Run("edge case for skipping Secondary Panel", func(t *testing.T) {
		m := initialModel()
		m.focusedPanel = MainPanel

		m.prevPanel()
		m.nextPanel()

		assertPanel(t, m.focusedPanel, MainPanel)
	})
}

// TestModel_Update tests the main update logic for key presses.
func TestModel_Update(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name          string
		initialPanel  Panel
		key           string
		expectedPanel Panel
	}{
		{
			name:          "Focus Next with Tab",
			initialPanel:  StatusPanel,
			key:           "tab",
			expectedPanel: FilesPanel,
		},
		{
			name:          "Focus Previous with Shift+Tab",
			initialPanel:  FilesPanel,
			key:           "shift+tab",
			expectedPanel: StatusPanel,
		},
		{
			name:          "Direct Focus with '0'",
			initialPanel:  StashPanel,
			key:           "0",
			expectedPanel: MainPanel,
		},
		{
			name:          "Direct Focus with '1'",
			initialPanel:  MainPanel,
			key:           "1",
			expectedPanel: StatusPanel,
		},
		{
			name:          "Direct Focus with '2'",
			initialPanel:  MainPanel,
			key:           "2",
			expectedPanel: FilesPanel,
		},
		{
			name:          "Direct Focus with '3'",
			initialPanel:  MainPanel,
			key:           "3",
			expectedPanel: BranchesPanel,
		},
		{
			name:          "Direct Focus with '4'",
			initialPanel:  MainPanel,
			key:           "4",
			expectedPanel: CommitsPanel,
		},
		{
			name:          "Direct Focus with '5'",
			initialPanel:  MainPanel,
			key:           "5",
			expectedPanel: StashPanel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := initialModel()
			m.focusedPanel = tc.initialPanel

			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune(tc.key),
			}
			if tc.key == "tab" {
				keyMsg.Type = tea.KeyTab
			}
			if tc.key == "shift+tab" {
				keyMsg.Type = tea.KeyShiftTab
			}

			updatedModel, _ := m.Update(keyMsg)
			newModel := updatedModel.(Model)

			if newModel.focusedPanel != tc.expectedPanel {
				t.Errorf("Update() with key '%s' failed: expected panel %v, got %v", tc.key, tc.expectedPanel, newModel.focusedPanel)
			}
		})
	}
}

func assertPanel(t testing.TB, got, want Panel) {
	if got != want {
		t.Errorf("nextPanel() failed: got %v, want %v", got, want)
	}
}
