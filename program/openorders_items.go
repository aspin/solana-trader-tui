package program

import (
	"fmt"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
)

type openOrdersItem struct {
	orderID       string
	side          string
	types         []pb.OrderType
	price         float64
	remainingSize float64
	clientOrderID string
}

func (i openOrdersItem) Title() string {
	return fmt.Sprintf("[%v] %v (%v)", i.side, i.orderID, i.clientOrderID)
}

func (i openOrdersItem) Description() string {
	return fmt.Sprintf("%v @ %v; types: %v", i.price, i.remainingSize, i.types)
}

func (i openOrdersItem) FilterValue() string {
	return i.orderID
}

func newOpenOrdersItem(order *pb.Order) openOrdersItem {
	return openOrdersItem{
		orderID:       order.OrderID,
		side:          order.Side.String(),
		types:         order.Types,
		price:         order.Price,
		remainingSize: order.RemainingSize,
		clientOrderID: order.ClientOrderID,
	}
}
