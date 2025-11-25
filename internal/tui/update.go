package tui

import (
	"github.com/cfstout/pr-watchtower/internal/gh"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.store.Close()
			return m, tea.Quit
		case "tab":
			m.activeList = (m.activeList + 1) % 2
			m.cursor = 0
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			var maxLen int
			if m.activeList == 0 {
				maxLen = len(m.needsReview)
			} else {
				maxLen = len(m.myPRs)
			}
			if m.cursor < maxLen-1 {
				m.cursor++
			}
		case "r":
			m.loading = true
			return m, fetchPRsCmd(m.cfg)
		case "a":
			// Trigger automation
			var selectedPR gh.PR
			if m.activeList == 0 {
				if len(m.needsReview) > 0 {
					selectedPR = m.needsReview[m.cursor]
				}
			} else {
				if len(m.myPRs) > 0 {
					selectedPR = m.myPRs[m.cursor]
				}
			}
			if selectedPR.Number != 0 {
				return m, func() tea.Msg {
					err := gh.TriggerWorkflow(selectedPR.Number)
					if err != nil {
						return errMsg{err}
					}
					return nil // Or a success message
				}
			}
		}

	case tickMsg:
		return m, tea.Batch(fetchPRsCmd(m.cfg), tickCmd(m.cfg))

	case prsLoadedMsg:
		m.loading = false
		m.needsReview = msg.incoming
		m.myPRs = msg.outgoing
		// Check status for each PR
		// Note: This is a side effect in Update, which is fine for now but could be async
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
