package program

import (
	"context"
	"github.com/aspin/solana-trader-tui/component/listquery"
	"github.com/aspin/solana-trader-tui/store"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	maxWidth = 80
)

type openOrdersModel struct {
	appStore *store.App
	dispatch StageDispatcher

	listquery listquery.Model
}

type openOrdersMsg struct {
	openOrders []*pb.Order
}

func newOpenOrdersModel(appStore *store.App) StageModel {
	marketInput := textinput.New()
	marketInput.Placeholder = "Market Name (e.g. SOL/USDC) or Public Key"
	marketInput.Focus()
	marketInput.PromptStyle = focusedStyle

	lq := listquery.New([]textinput.Model{marketInput}, spinner.Points, list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0), nil)

	m := &openOrdersModel{
		appStore:  appStore,
		listquery: lq,
	}
	m.listquery.SetQuery(m.fetchOrders)
	return m
}

func (m *openOrdersModel) Init(dispatch StageDispatcher) tea.Cmd {
	m.dispatch = dispatch
	return m.listquery.Init(m.appStore.UI.WindowWidth, m.appStore.UI.WindowHeight)
}

func (m *openOrdersModel) Update(msg tea.Msg) (Stage, StageModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		exit bool
	)

	m.listquery, cmd, exit = m.listquery.Update(msg)
	if exit {
		return StageMenu, m, nil
	}
	return StageOpenOrders, m, cmd
}

func (m *openOrdersModel) fetchOrders(vs []string) {
	market := vs[0]
	openOrders, err := m.appStore.Provider.GetOpenOrders(context.Background(), market, "", m.appStore.Settings.OpenOrdersAddress.String(), m.appStore.Settings.Project)
	if err != nil {
		m.dispatch(listquery.ErrorMsg{Err: err})
		return
	}

	items := make([]list.Item, 0)
	for _, order := range openOrders.Orders {
		items = append(items, newOpenOrdersItem(order))
	}
	m.dispatch(listquery.ResultMsg{Items: items})
}

func (m openOrdersModel) View() string {
	return m.listquery.View()
}
