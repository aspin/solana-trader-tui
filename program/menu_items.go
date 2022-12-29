package program

import "github.com/charmbracelet/bubbles/list"

type menuItem struct {
	title string
	desc  string
}

func (i menuItem) Title() string {
	return i.title
}

func (i menuItem) Description() string {
	return i.desc
}

func (i menuItem) FilterValue() string {
	return i.title
}

var defaultMenuItems = []list.Item{
	menuItem{
		title: "Settings",
		desc:  "Set app details such as private/public key, auth header, etc.",
	},
	menuItem{
		title: "Open Orders",
		desc:  "View your unfilled open orders in a dex market",
	},
	menuItem{
		title: "Orderbook",
		desc:  "View all asks and bids in a dex market",
	},
	menuItem{
		title: "Stream Orderbook",
		desc:  "View stream of orderbook updates in a dex market",
	},
}
