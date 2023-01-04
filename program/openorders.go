package program

import (
	"context"
	"fmt"
	"github.com/aspin/solana-trader-tui/store"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

const (
	maxWidth = 80
)

type ooState int

const (
	ooInput ooState = iota
	ooLoading
	ooShow
)

type openOrdersModel struct {
	appStore *store.App
	dispatch StageDispatcher

	marketInput textinput.Model
	progress    progress.Model
	list        list.Model

	state ooState
	err   error
}

type openOrdersMsg struct {
	openOrders []*pb.Order
}

type tickMsg struct{}

func newOpenOrdersModel(appStore *store.App) StageModel {
	marketInput := textinput.New()
	marketInput.Placeholder = "Market Name (e.g. SOL/USDC) or Public Key"
	marketInput.Focus()
	marketInput.PromptStyle = focusedStyle

	m := &openOrdersModel{
		appStore:    appStore,
		marketInput: marketInput,
		progress:    progress.New(progress.WithDefaultGradient()),
		list:        list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		state:       ooInput,
	}

	m.list.Title = "Open Orders"

	// FIXME: title is current sticky?
	m.list.SetShowTitle(false)
	return m
}

func (m *openOrdersModel) Init(dispatch StageDispatcher) tea.Cmd {
	m.dispatch = dispatch
	m.state = ooInput
	m.marketInput.SetValue("")
	m.progress.SetPercent(0)
	m.err = nil
	m.setSize()
	return textinput.Blink
}

func (m *openOrdersModel) setSize() {
	h, v := listStyle.GetFrameSize()
	m.list.SetSize(m.appStore.UI.WindowWidth-h, m.appStore.UI.WindowHeight-v)

	m.progress.Width = m.appStore.UI.WindowWidth - 4
	if m.progress.Width > maxWidth {
		m.progress.Width = maxWidth
	}
}

func (m *openOrdersModel) Update(msg tea.Msg) (Stage, StageModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize()
	case tea.KeyMsg:
		switch m.state {
		case ooInput:
			if msg.Type == tea.KeyEnter {
				m.state = ooLoading
				m.fetchOrders()
			}
		case ooShow:
			if key.Matches(msg, m.list.KeyMap.Quit) && !m.list.IsFiltered() {
				return StageMenu, m, nil
			}
		}
	case tickMsg:
		switch m.state {
		case ooLoading:
			cmd = m.progress.IncrPercent(0.2)
			return StageOpenOrders, m, cmd
		}
	case openOrdersMsg:
		items := make([]list.Item, 0)
		for _, order := range msg.openOrders {
			items = append(items, newOpenOrdersItem(order))
		}
		m.list.SetItems(items)
		m.state = ooShow
	case progress.FrameMsg:
		switch m.state {
		case ooLoading:
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return StageOpenOrders, m, cmd
		}
		if m.state == ooLoading {
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return StageOpenOrders, m, cmd
		}
	}

	switch m.state {
	case ooInput:
		m.marketInput, cmd = m.marketInput.Update(msg)
		return StageOpenOrders, m, cmd
	case ooLoading:
		return StageOpenOrders, m, nil
	case ooShow:
		m.list, cmd = m.list.Update(msg)
		return StageOpenOrders, m, cmd
	default:
		panic(fmt.Errorf("open orders reached unknown state: %v", m.state))
	}
}

func (m *openOrdersModel) fetchOrders() {
	go func() {
		openOrders, err := m.appStore.Provider.GetOpenOrders(context.Background(), m.marketInput.Value(), "", m.appStore.Settings.OpenOrdersAddress.String(), m.appStore.Settings.Project)
		if err != nil {
			m.err = err
			m.state = ooInput
			m.progress.SetPercent(0)
			return
		}
		m.dispatch(openOrdersMsg{openOrders: openOrders.Orders})
	}()
	go func() {
		// completely artificial loading bar
		for m.state == ooLoading {
			time.Sleep(200 * time.Millisecond)
			m.dispatch(tickMsg{})
		}
	}()
}

func (m openOrdersModel) View() string {
	var b strings.Builder

	switch m.state {
	case ooInput:
		b.WriteString(m.marketInput.View())
		b.WriteRune('\n')

		if m.err != nil {
			b.WriteString(errorStyle.Render(m.err.Error()))
			b.WriteRune('\n')
		}
	case ooLoading:
		_, _ = fmt.Fprintf(&b, "Loading open orders for SOL/USDC (%v) for %v...\n", m.appStore.Settings.Project, m.appStore.Settings.PublicKey)
		b.WriteString(m.progress.View())
		b.WriteString("\n\n")
	case ooShow:
		b.WriteString(listStyle.Render(m.list.View()))
	}

	return b.String()
}
