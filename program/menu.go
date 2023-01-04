package program

import (
	"github.com/aspin/solana-trader-tui/store"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type menuModel struct {
	appStore *store.App
	list     list.Model
}

func newMenuModel(appStore *store.App) StageModel {
	m := &menuModel{
		appStore: appStore,
		list:     list.New(defaultMenuItems, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Menu"

	// FIXME: title is currently sticky?
	m.list.SetShowTitle(false)
	return m
}

func (m *menuModel) Init(dispatch StageDispatcher) tea.Cmd {
	m.setSize()
	return nil
}

func (m *menuModel) setSize() {
	h, v := listStyle.GetFrameSize()
	m.list.SetSize(m.appStore.UI.WindowWidth-h, m.appStore.UI.WindowHeight-v)
}

func (m *menuModel) Update(msg tea.Msg) (Stage, StageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize()
	case tea.KeyMsg:
		switch msg.Type {
		// transition to other stage
		case tea.KeyEnter, tea.KeySpace:
			listIndex := m.list.Index()
			stage := m.list.Items()[listIndex].(menuItem).stage
			return stage, m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return StageMenu, m, cmd
}

func (m *menuModel) View() string {
	return listStyle.Render(m.list.View())
}
