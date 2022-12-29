package program

import (
	"context"
	"fmt"
	"github.com/aspin/solana-trader-tui/store"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
	"time"
)

const (
	maxWidth = 80
)

type openOrdersModel struct {
	appStore *store.App
	list     list.Model
	loading  bool
	progress progress.Model
	err      error // TODO: unused right now
}

type openOrdersMsg struct {
	openOrders []*pb.Order
}

type tickMsg struct{}

func newOpenOrdersModel(appStore *store.App) StageModel {
	m := &openOrdersModel{
		appStore: appStore,
		list:     list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		progress: progress.New(progress.WithDefaultGradient()),
		loading:  true,
	}

	m.list.Title = "Open Orders"

	// FIXME: title is current sticky?
	m.list.SetShowTitle(false)
	return m
}

func (m *openOrdersModel) Init(dispatch StageDispatcher) tea.Cmd {
	m.loading = true

	go func() {
		openOrders, err := m.appStore.Provider.GetOpenOrders(context.Background(), "SOL/USDC", "", m.appStore.Settings.OpenOrdersAddress.String(), m.appStore.Settings.Project)
		time.Sleep(time.Second)
		m.loading = false
		if err != nil {
			log.Printf("could not fetch orders: %v", err)
			return
		}
		dispatch(openOrdersMsg{openOrders: openOrders.Orders})
	}()
	go func() {
		for m.loading {
			time.Sleep(200 * time.Millisecond)
			dispatch(tickMsg{})
		}
	}()
	m.setSize()
	return nil
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize()
	case tickMsg:
		if m.loading {
			cmd := m.progress.IncrPercent(0.2)
			return StageOpenOrders, m, cmd
		}
	case openOrdersMsg:
		items := make([]list.Item, 0)
		for _, order := range msg.openOrders {
			items = append(items, newOpenOrdersItem(order))
		}
		m.list.SetItems(items)
	case progress.FrameMsg:
		if m.loading {
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return StageOpenOrders, m, cmd
		}
	}

	// TODO: a bit not elegant
	if m.loading {
		return StageOpenOrders, m, nil
	} else {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return StageOpenOrders, m, cmd
	}
}

func (m openOrdersModel) View() string {
	if m.loading {
		var b strings.Builder
		_, _ = fmt.Fprintf(&b, "Loading open orders for SOL/USDC (%v) for %v...\n", m.appStore.Settings.Project, m.appStore.Settings.PublicKey)
		b.WriteString(m.progress.View())
		b.WriteString("\n\n")

		return b.String()
	} else {
		return listStyle.Render(m.list.View())
	}
}
