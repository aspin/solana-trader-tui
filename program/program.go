package program

import (
	"github.com/aspin/solana-trader-tui/store"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

type appModel struct {
	stage    Stage
	models   map[Stage]StageModel
	store    *store.App
	dispatch StageDispatcher
}

func New() *tea.Program {
	s := &store.App{}

	// TODO: add loading from config file
	initialStage := StageMenu
	if s.NeedsInit() {
		initialStage = StageSettings
	}

	m := &appModel{
		stage: initialStage,
		store: s,
	}

	models := map[Stage]StageModel{
		StageSettings: newSettingsModel(m.store),
		StageMenu:     newMenuModel(m.store),
		StageError:    newErrorModel(m.store),
	}
	m.models = models

	p := tea.NewProgram(m, tea.WithAltScreen())
	m.dispatch = p.Send
	return p
}

func (m appModel) Init() tea.Cmd {
	return m.models[m.stage].Init(m.dispatch)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.store.UI.WindowWidth = msg.Width
		m.store.UI.WindowHeight = msg.Height
	}

	// process current stage updates
	model, ok := m.models[m.stage]
	if !ok {
		log.Printf("error[update]: could not find model for stage %v", m.stage)
		return m, tea.Quit
	}

	nextStage, nextModel, nextCmd := model.Update(msg)
	if nextStage == StageExit {
		return m, tea.Quit
	}
	m.models[m.stage] = nextModel

	// no stage transition; continue
	if m.stage == nextStage {
		return m, nextCmd
	}

	// stage transition: move onto initializing next model
	nextModel, ok = m.models[nextStage]
	if !ok {
		log.Printf("error[update]: could not find model for next stage %v", m.stage)
		return m, tea.Quit
	}
	m.stage = nextStage
	return m, nextModel.Init(m.dispatch)
}

func (m appModel) View() string {
	var b strings.Builder
	b.WriteString("bloXroute Trader API\n\n")

	model, ok := m.models[m.stage]
	if !ok {
		log.Printf("error[view]: could not find model for stage %v", m.stage)
		return ""
	}

	b.WriteString(model.View())
	return b.String()
}
