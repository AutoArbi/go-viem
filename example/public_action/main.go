package main

import (
	"context"
	"fmt"
	"github.com/AutoArbi/go-viem/eth"
	log2 "github.com/ethereum/go-ethereum/log"
	"log"
	"time"

	"github.com/AutoArbi/go-viem/client"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// 创建 WebSocket 传输
	wsTransport, err := client.NewWebSocketTransport("wss://eth-mainnet.g.alchemy.com/v2/LE5jeRzEtZ890kR9v9wG4TXQ6itwoDZu")
	if err != nil {
		log.Fatalf("Failed to create WebSocket transport: %v", err)
	}

	// 创建 JSON-RPC 客户端
	cli, err := client.NewClient(client.WithTransport(wsTransport))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 创建公共客户端
	publicClient := &eth.Client{Client: cli}

	// 设置上下文和超时
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// eth_blockNumber
	blockNumber, err := publicClient.GetBlockNumber(ctx)
	if err != nil {
		log.Fatalf("Failed to get block number: %v", err)
	}
	fmt.Printf("Block Number: %s\n", blockNumber)

	blockByNumber, err := publicClient.GetBlockByNumber(ctx, blockNumber, true)
	if err != nil {
		log.Fatalf("Failed to get block by number: %v", err)
	}
	fmt.Printf("Block By Number: %s\n", blockByNumber)

	log2.Trace("Block By Number", blockByNumber)

	// eth_getChainId
	chainId, err := publicClient.GetChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}
	fmt.Printf("Chain ID: %s\n", chainId.String())

	// 获取账户余额
	address := common.HexToAddress("0x7031576A278AfFcdEcf0d6E9673E51261C0BFFF8")
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
