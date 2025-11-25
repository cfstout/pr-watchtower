package tui

import (
	"os/exec"
	"time"

	"github.com/cfstout/pr-watchtower/internal/gh"
	"github.com/cfstout/pr-watchtower/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Helper to get current visible list
		var currentList []gh.PR
		if m.activeList == 0 {
			currentList = getVisiblePRs(m.store, m.needsReview, m.showHidden)
		} else {
			currentList = getVisiblePRs(m.store, m.myPRs, m.showHidden)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.store.Close()
			return m, tea.Quit
		case "tab":
			// Check if other list has visible items?
			// User logic: "Only allow switching if both lists have items"
			// We should probably allow switching if the other list exists, even if empty?
			// But original logic checked len > 0.
			// Let's stick to original logic but maybe check visible len?
			// For now, keep original logic for switching availability, but reset cursor.
			if len(m.needsReview) > 0 && len(m.myPRs) > 0 {
				m.activeList = (m.activeList + 1) % 2
				m.cursor = 0
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(currentList)-1 {
				m.cursor++
			}
		case "r":
			m.loading = true
			return m, fetchPRsCmd(m.cfg)
		case "h":
			if len(currentList) > 0 && m.cursor < len(currentList) {
				selectedPR := currentList[m.cursor]
				isHidden, _ := m.store.IsHidden(selectedPR.Number)
				m.store.SetHidden(selectedPR.Number, !isHidden)

				// If we just hid it, and showHidden is false, it will disappear.
				// We need to clamp cursor if it's now out of bounds.
				// Actually, if it disappears, the list shrinks by 1.
				// If we were at the last item, we need to decrement cursor.
				if !m.showHidden && !isHidden { // We just hid it
					// List length will decrease by 1
					newLen := len(currentList) - 1
					if m.cursor >= newLen && newLen > 0 {
						m.cursor = newLen - 1
					} else if newLen == 0 {
						m.cursor = 0
					}
				}
			}
		case "H":
			m.showHidden = !m.showHidden
			// Re-calculate visible list to clamp cursor
			var newList []gh.PR
			if m.activeList == 0 {
				newList = getVisiblePRs(m.store, m.needsReview, m.showHidden)
			} else {
				newList = getVisiblePRs(m.store, m.myPRs, m.showHidden)
			}
			if m.cursor >= len(newList) {
				if len(newList) > 0 {
					m.cursor = len(newList) - 1
				} else {
					m.cursor = 0
				}
			}
		case "enter":
			if len(currentList) > 0 && m.cursor < len(currentList) {
				selectedPR := currentList[m.cursor]
				return m, func() tea.Msg {
					// Use 'open' command on macOS
					exec.Command("open", selectedPR.Url).Run()
					return nil
				}
			}
		case "m":
			if len(currentList) > 0 && m.cursor < len(currentList) {
				selectedPR := currentList[m.cursor]
				return m, func() tea.Msg {
					err := gh.MergePR(selectedPR.Number)
					if err != nil {
						return errMsg{err}
					}
					return tickMsg(time.Now())
				}
			}
		case "a":
			if len(currentList) > 0 && m.cursor < len(currentList) {
				selectedPR := currentList[m.cursor]
				return m, func() tea.Msg {
					err := gh.TriggerWorkflow(selectedPR.Number)
					if err != nil {
						return errMsg{err}
					}
					return nil
				}
			}
		}

	case tickMsg:
		return m, tea.Batch(fetchPRsCmd(m.cfg), tickCmd(m.cfg))

	case prsLoadedMsg:
		m.loading = false
		m.needsReview = msg.incoming
		m.myPRs = msg.outgoing

		// Auto-focus logic
		if m.activeList == 0 && len(m.needsReview) == 0 && len(m.myPRs) > 0 {
			m.activeList = 1
			m.cursor = 0
		} else if m.activeList == 1 && len(m.myPRs) == 0 && len(m.needsReview) > 0 {
			m.activeList = 0
			m.cursor = 0
		}

		// Update DB status
		for _, pr := range m.needsReview {
			m.store.CheckUpdateStatus(pr.Number, pr.UpdatedAt)
		}
		for _, pr := range m.myPRs {
			m.store.CheckUpdateStatus(pr.Number, pr.UpdatedAt)
		}

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func getVisiblePRs(s *store.Store, prs []gh.PR, showHidden bool) []gh.PR {
	if showHidden {
		return prs
	}
	var visible []gh.PR
	for _, pr := range prs {
		hidden, _ := s.IsHidden(pr.Number)
		if !hidden {
			visible = append(visible, pr)
		}
	}
	return visible
}
