package program

import (
	"github.com/aspin/solana-trader-tui/store"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("197"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
)

type errorModel struct {
	appStore *store.App
}

func newErrorModel(appStore *store.App) StageModel {
	return errorModel{appStore: appStore}
}

func (m errorModel) Init(dispatch StageDispatcher) tea.Cmd {
	return nil
}

func (m errorModel) Update(msg tea.Msg) (Stage, StageModel, tea.Cmd) {
	return StageError, m, nil
}

func (m errorModel) View() string {
	var b strings.Builder

	b.WriteString(errorStyle.Render("Encountered fatal error:"))
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(noStyle.Render(m.appStore.Err.Error()))
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(helpStyle.Render("(ctrl+c or esc to exit)"))
	return b.String()
}
