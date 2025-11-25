package tui

import (
	"fmt"

	"github.com/cfstout/pr-watchtower/internal/gh"
	"github.com/cfstout/pr-watchtower/internal/store"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	listHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(lipgloss.Color("170")).
				SetString("> ")

	statusNewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green
			SetString("[NEW]")

	statusUpdatedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("208")). // Orange
				SetString("[UPDATED]")

	statusSeenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Grey
			SetString("[SEEN]")
)

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	s := titleStyle.Render("PR Watchtower") + "\n\n"

	if m.loading {
		s += "Loading PRs...\n"
	} else {
		s += m.renderList("Needs Review", m.needsReview, m.activeList == 0)
		s += "\n"
		s += m.renderList("My PRs", m.myPRs, m.activeList == 1)
	}

	s += "\n" + helpStyle.Render("r: refresh • a: agent fix • tab: switch list • q: quit") + "\n"
	return s
}

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	MarginTop(1)

func (m Model) renderList(title string, prs []gh.PR, active bool) string {
	s := listHeaderStyle.Render(title) + "\n"
	if len(prs) == 0 {
		s += itemStyle.Render("(No PRs)") + "\n"
		return s
	}

	for i, pr := range prs {
		cursor := "  "
		style := itemStyle
		if active && m.cursor == i {
			cursor = "> "
			style = selectedItemStyle
		}

		status, _ := m.store.CheckUpdateStatus(pr.Number, pr.UpdatedAt)
		statusBadge := ""
		switch status {
		case store.StatusNew:
			statusBadge = statusNewStyle.Render()
		case store.StatusUpdated:
			statusBadge = statusUpdatedStyle.Render()
		default:
			statusBadge = statusSeenStyle.Render()
		}

		line := fmt.Sprintf("%s %s #%d %s", cursor, statusBadge, pr.Number, pr.Title)
		s += style.Render(line) + "\n"
	}
	return s
}
