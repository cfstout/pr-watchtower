package tui

import (
	"time"

	"github.com/cfstout/pr-watchtower/internal/config"
	"github.com/cfstout/pr-watchtower/internal/gh"
	"github.com/cfstout/pr-watchtower/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type Model struct {
	cfg         *config.Config
	store       *store.Store
	needsReview []gh.PR
	myPRs       []gh.PR
	cursor      int
	activeList  int // 0 for Needs Review, 1 for My PRs
	loading     bool
	err         error
	width       int
	height      int
}

func InitialModel(cfg *config.Config) Model {
	s, err := store.NewStore("")
	if err != nil {
		// In a real app, handle this better
		panic(err)
	}
	return Model{
		cfg:        cfg,
		store:      s,
		loading:    true,
		activeList: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchPRsCmd(m.cfg),
		tickCmd(m.cfg),
	)
}

func tickCmd(cfg *config.Config) tea.Cmd {
	return tea.Tick(cfg.GitHub.RefreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchPRsCmd(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		// Fetch "Needs Review"
		incoming, err := gh.FetchPRs(cfg.GitHub.Queries.NeedsReview)
		if err != nil {
			return errMsg{err}
		}

		// Fetch "My PRs"
		outgoing, err := gh.FetchPRs(cfg.GitHub.Queries.MyPRs)
		if err != nil {
			return errMsg{err}
		}

		return prsLoadedMsg{incoming: incoming, outgoing: outgoing}
	}
}

type prsLoadedMsg struct {
	incoming []gh.PR
	outgoing []gh.PR
}

type errMsg struct{ err error }
