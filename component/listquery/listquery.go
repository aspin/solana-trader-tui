package listquery

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("197"))
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	listStyle    = lipgloss.NewStyle().Margin(1, 2)
	noStyle      = lipgloss.NewStyle()
)

type viewState int

const (
	vsInput viewState = iota
	vsLoading
	vsShow
)

type queryFn func([]string)

type ResultMsg struct {
	Items []list.Item
}

type ErrorMsg struct {
	Err error
}

type Model struct {
	focusIndex int
	inputs     []textinput.Model

	spinner spinner.Model
	list    list.Model
	query   queryFn

	state viewState
	err   error
}

func New(inputs []textinput.Model, spinnerType spinner.Spinner, l list.Model, query queryFn) Model {
	l.SetShowTitle(false)
	return Model{
		inputs:  inputs,
		spinner: spinner.New(spinner.WithSpinner(spinnerType)),
		list:    l,
		state:   vsInput,
		query:   query,
	}
}

func (m *Model) SetQuery(query queryFn) {
	m.query = query
}

func (m *Model) Init(width, height int) tea.Cmd {
	m.focusIndex = 0
	m.state = vsInput
	for _, input := range m.inputs {
		input.SetValue("")
	}
	m.focusInputs()

	m.err = nil
	m.setSize(width, height)
	return textinput.Blink
}

func (m *Model) setSize(width, height int) {
	h, v := listStyle.GetFrameSize()
	m.list.SetSize(width-h, height-v)
}

func (m *Model) focusInputs() tea.Cmd {
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
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd, bool) {
	var cmd tea.Cmd

	// applies to all message types
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
	}

	switch m.state {
	case vsInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch k := msg.Type; k {
			case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
				if msg.Type == tea.KeyEnter && m.focusIndex == len(m.inputs) {
					if err := m.validateInputs(); err == nil {
						m.state = vsLoading
						go m.query(m.inputValues())
						return m, m.spinner.Tick, false
					} else {
						m.err = err
					}
				}

				m.updateCursor(k)
				cmd = m.focusInputs()
				return m, cmd, false
			}
		}
		m.updateInputs(msg)
		return m, cmd, false
	case vsLoading:
		switch msg := msg.(type) {
		case spinner.TickMsg:
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd, false
		case ResultMsg:
			m.list.SetItems(msg.Items)
			m.state = vsShow
		case ErrorMsg:
			m.err = msg.Err
			m.state = vsInput
			cmd = textinput.Blink
		}
		return m, cmd, false
	case vsShow:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.list.KeyMap.Quit) && !m.list.IsFiltered() {
				return m, nil, true
			}
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd, false
	default:
		panic(fmt.Errorf("list query reached unknown state: %v", m.state))
	}
}

func (m Model) validateInputs() error {
	for _, input := range m.inputs {
		if input.Err != nil {
			return input.Err
		}
	}
	return nil
}

func (m Model) inputValues() []string {
	s := make([]string, 0, len(m.inputs))
	for _, input := range m.inputs {
		s = append(s, input.Value())
	}
	return s
}

func (m *Model) updateCursor(k tea.KeyType) {
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
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	switch s := m.state; s {
	case vsInput, vsLoading:
		for _, input := range m.inputs {
			b.WriteString(input.View())
			b.WriteRune('\n')
		}

		button := "\n[ Submit ]\n"
		if m.focusIndex == len(m.inputs) {
			button = focusedStyle.Render(button)
		}
		b.WriteString(button)
		b.WriteRune('\n')

		if m.err != nil {
			b.WriteString(errorStyle.Render(m.err.Error()))
			b.WriteRune('\n')
		}

		if s == vsLoading {
			b.WriteString(m.spinner.View())
		}
	case vsShow:
		b.WriteString(listStyle.Render(m.list.View()))
	}

	return b.String()
}
