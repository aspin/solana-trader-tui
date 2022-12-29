package program

import (
	"github.com/aspin/solana-trader-tui/store"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"strings"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	noStyle      = lipgloss.NewStyle()
)

type settingsModel struct {
	inputs     []textinput.Model
	focusIndex int

	appStore *store.App
}

func newSettingsModel(appStore *store.App) StageModel {
	m := settingsModel{
		inputs:   make([]textinput.Model, 5),
		appStore: appStore,
	}

	for i := range m.inputs {
		t := textinput.New()
		switch i {
		case 0:
			t.Placeholder = "bloXroute Auth Header"
			t.Focus()
			t.PromptStyle = focusedStyle
		case 1:
			t.Placeholder = "Private Key"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		case 2:
			t.Placeholder = "Public Key"
		case 3:
			t.Placeholder = "Open Orders Address"
			// TODO: validate is address
		case 4:
			t.Placeholder = "Project"
			// TODO: validate project is project
		}

		m.inputs[i] = t
	}
	return m
}

func (m settingsModel) Init(dispatch StageDispatcher) tea.Cmd {
	return textinput.Blink
}

func (m settingsModel) Update(msg tea.Msg) (Stage, StageModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch k := msg.Type; k {
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			if k == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				err := m.submit()
				if err != nil {
					log.Printf("could not update settings: %v", err)
				}
				return StageMenu, m, nil
			}

			if k == tea.KeyUp || k == tea.KeyShiftTab {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					continue
				}

				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
			}
			return StageSettings, m, tea.Batch(cmds...)
		}
	}

	cmd = m.updateInputs(msg)
	return StageSettings, m, cmd
}

func (m *settingsModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m settingsModel) submit() error {
	for _, input := range m.inputs {
		if input.Err != nil {
			return input.Err
		}
	}

	m.appStore.Settings.AuthHeader = m.inputs[0].Value()
	m.appStore.Settings.PrivateKey = m.inputs[1].Value()
	m.appStore.Settings.PublicKey = m.inputs[2].Value()
	m.appStore.Settings.OpenOrdersAddress = m.inputs[3].Value()
	return nil
}

func (m settingsModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	button := "\n[ Submit ]\n"
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Render(button)
	}
	b.WriteString(button)
	return b.String()
}
