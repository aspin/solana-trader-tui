package program

import tea "github.com/charmbracelet/bubbletea"

type Stage int

// StageDispatcher sends messages to the global context to be processed
type StageDispatcher func(msg tea.Msg)

type StageModel interface {
	Init(dispatch StageDispatcher) tea.Cmd

	Update(msg tea.Msg) (Stage, StageModel, tea.Cmd)

	View() string
}

var (
	StageExit       Stage = 0
	StageMenu       Stage = 1
	StageSettings   Stage = 2
	StageOpenOrders Stage = 3
	StageView       Stage = 4
	StageError      Stage = 5
)
