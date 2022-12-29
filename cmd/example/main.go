package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
)

// An example of what a fairly plain "interactive" application would look like.

func main() {

	var (
		authHeader        string
		market            string
		openOrdersAddress string
		err               error
	)

	fmt.Println("Getting open orders for account!")
	fmt.Println("")

	fmt.Println("Enter your auth header: ")
	_, err = fmt.Scanln(&authHeader)
	if err != nil {
		panic(err)
	}

	fmt.Println("Enter the market you want to check: ")
	_, err = fmt.Scanln(&market)
	if err != nil {
		panic(err)
	}

	fmt.Println("Enter your open orders address: ")
	_, err = fmt.Scanln(&openOrdersAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	fmt.Printf("Loading...")

	opts := provider.RPCOpts{
		Endpoint:   provider.MainnetGRPC,
		UseTLS:     true,
		AuthHeader: authHeader,
	}
	g, err := provider.NewGRPCClientWithOpts(opts)
	if err != nil {
		panic(err)
	}

	orders, err := g.GetOpenOrders(context.Background(), market, "", openOrdersAddress, pb.Project_P_SERUM)
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	fmt.Printf("Orders (%v): \n", len(orders.Orders))
	for _, order := range orders.Orders {
		fmt.Printf("  [%v] Order %v: %v@%v\n", order.Side, order.OrderID, order.RemainingSize, order.Price)
	}
}
