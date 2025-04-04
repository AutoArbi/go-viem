package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AutoArbi/go-viem/action/public"
	"github.com/AutoArbi/go-viem/client"
	"github.com/AutoArbi/go-viem/transport"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// 创建 WebSocket 传输
	wsTransport, err := transport.NewWebSocketTransport("wss://example.com/rpc")
	if err != nil {
		log.Fatalf("Failed to create WebSocket transport: %v", err)
	}

	// 创建 JSON-RPC 客户端
	cli, err := client.NewClient(client.WithTransport(wsTransport))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 创建公共客户端
	publicClient := &public.Client{Client: cli}

	// 设置上下文和超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取账户余额
	address := common.HexToAddress("0xYourEthereumAddress")
	balance, err := publicClient.GetBalance(ctx, address, "latest")
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("Balance: %s\n", balance.String())

	// 获取交易计数
	txCount, err := publicClient.GetTransactionCount(ctx, address, "latest")
	if err != nil {
		log.Fatalf("Failed to get transaction count: %v", err)
	}
	fmt.Printf("Transaction Count: %d\n", txCount)
}
