package program

import (
	"errors"
	"fmt"
	"github.com/aspin/solana-trader-tui/store"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gagliardetto/solana-go"
	"strings"
	"time"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	noStyle      = lipgloss.NewStyle()
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
)

type settingsModel struct {
	err    error
	status string

	inputs     []textinput.Model
	focusIndex int

	appStore *store.App
	dispatch StageDispatcher
}

type statusMsg struct {
	status string
}

type statusErrMsg struct {
	err error
}

type statusDoneMsg struct{}

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
			t.SetValue(appStore.Settings.AuthHeader)
		case 1:
			t.Placeholder = "Private Key"
			t.SetValue(appStore.Settings.PrivateKey.String())
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		case 2:
			t.Placeholder = "Public Key"

			publicKey := appStore.Settings.PublicKey
			if !publicKey.IsZero() {
				t.SetValue(publicKey.String())
			}
		case 3:
			t.Placeholder = "Open Orders Address"

			openOrdersAddress := appStore.Settings.OpenOrdersAddress
			if !openOrdersAddress.IsZero() {
				t.SetValue(openOrdersAddress.String())
			}
		case 4:
			t.Placeholder = "Project"

			project := appStore.Settings.Project
			if project != pb.Project_P_UNKNOWN {
				t.SetValue(project.String())
			}
		}

		m.inputs[i] = t
	}

	m.Init(nil)
	return &m
}

func (m *settingsModel) Init(dispatch StageDispatcher) tea.Cmd {
	m.focusIndex = 0
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = focusedStyle
	m.dispatch = dispatch
	return textinput.Blink
}

func (m *settingsModel) Update(msg tea.Msg) (Stage, StageModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch k := msg.Type; k {
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			if k == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				err := m.submit()
				if err != nil {
					m.err = err
					break
				}
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
	case statusMsg:
		m.status = msg.status
	case statusErrMsg:
		m.err = msg.err
	case statusDoneMsg:
		return StageMenu, m, nil
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
	authHeader := m.inputs[0].Value()
	if authHeader == "" {
		return errors.New("auth header cannot be empty")
	}
	m.appStore.Settings.AuthHeader = authHeader

	privateKeyStr := m.inputs[1].Value()
	if privateKeyStr != "" {
		privateKey, err := solana.PrivateKeyFromBase58(privateKeyStr)
		if err != nil {
			return fmt.Errorf("invalid private key: %w", err)
		}
		m.appStore.Settings.PrivateKey = privateKey
	} else {
		m.appStore.Settings.PrivateKey = nil
	}

	publicKeyStr := m.inputs[2].Value()
	publicKey, err := solana.PublicKeyFromBase58(publicKeyStr)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	m.appStore.Settings.PublicKey = publicKey

	openOrdersAddressStr := m.inputs[3].Value()
	openOrdersAddress, err := solana.PublicKeyFromBase58(openOrdersAddressStr)
	if err != nil {
		return fmt.Errorf("invalid open orders address key: %w", err)
	}
	m.appStore.Settings.OpenOrdersAddress = openOrdersAddress

	projectStr := m.inputs[4].Value()
	projectInt, ok := pb.Project_value[projectStr]
	project := pb.Project(projectInt)
	if !ok || project == pb.Project_P_UNKNOWN {
		return fmt.Errorf("invalid project value: %v", projectStr)
	}
	m.appStore.Settings.Project = project

	go func() {
		m.dispatch(statusMsg{status: "connecting..."})
		err = m.appStore.Connect()
		if err != nil {
			m.dispatch(statusErrMsg{err: err})
			return
		}
		m.dispatch(statusMsg{status: "connected!"})
		time.Sleep(500 * time.Millisecond)
		m.dispatch(statusDoneMsg{})
	}()
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
	b.WriteRune('\n')

	if m.status != "" {
		b.WriteString(statusStyle.Render(m.status))
		b.WriteRune('\n')
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()))
		b.WriteRune('\n')
	}
	return b.String()
}
