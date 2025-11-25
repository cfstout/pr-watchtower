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

	statusSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")). // Green
				SetString("✔")

	statusPendingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("220")). // Yellow
				SetString("●")

	statusFailureStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")). // Red
				SetString("✖")

	statusConflictStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("208")). // Orange
				SetString("!")

	hiddenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			SetString("[HIDDEN]")
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

	helpText := "r: refresh • a: agent fix • h: hide • H: show hidden • enter: open • m: merge • q: quit"
	if len(m.needsReview) > 0 && len(m.myPRs) > 0 {
		helpText = "tab: switch • " + helpText
	}
	s += "\n" + helpStyle.Render(helpText) + "\n"
	return s
}

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	MarginTop(1)

func (m Model) renderList(title string, prs []gh.PR, active bool) string {
	s := listHeaderStyle.Render(title) + "\n"

	// Filter visible PRs
	// We duplicate logic here or make getVisiblePRs a method of Model or exported function.
	// Since it's in the same package (tui), we can use it if it's defined in update.go (which is package tui).
	// Go allows calling functions defined in other files of the same package.
	visiblePRs := getVisiblePRs(m.store, prs, m.showHidden)

	if len(visiblePRs) == 0 {
		s += itemStyle.Render("(No PRs)") + "\n"
		return s
	}

	for i, pr := range visiblePRs {
		cursor := "  "
		style := itemStyle

		// Cursor logic: m.cursor is now index into visiblePRs
		if active && m.cursor == i {
			cursor = "> "
			style = selectedItemStyle
		}

		status, _ := m.store.CheckUpdateStatus(pr.Number, pr.UpdatedAt)
		statusBadge := ""

		// Build Status / Mergeability
		buildBadge := ""
		if pr.Mergeable == "CONFLICTING" {
			buildBadge = statusConflictStyle.Render()
		} else if pr.StatusCheckRollup.State == "FAILURE" {
			buildBadge = statusFailureStyle.Render()
		} else if pr.StatusCheckRollup.State == "PENDING" {
			buildBadge = statusPendingStyle.Render()
		} else if pr.StatusCheckRollup.State == "SUCCESS" {
			buildBadge = statusSuccessStyle.Render()
		}

		switch status {
		case store.StatusNew:
			statusBadge = statusNewStyle.Render()
		case store.StatusUpdated:
			statusBadge = statusUpdatedStyle.Render()
		default:
			statusBadge = statusSeenStyle.Render()
		}

		hiddenBadge := ""
		isHidden, _ := m.store.IsHidden(pr.Number)
		if isHidden {
			hiddenBadge = hiddenStyle.Render() + " "
		}

		line := fmt.Sprintf("%s %s %s%s #%d %s", cursor, statusBadge, hiddenBadge, buildBadge, pr.Number, pr.Title)
		s += style.Render(line) + "\n"
	}
	return s
}
