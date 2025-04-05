package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AutoArbi/go-viem/client"
)

func main() {
	httpTransport, err := client.NewWebSocketTransport("wss://example.com/rpc")
	if err != nil {
		log.Fatalf("Failed to create WebScoket transport: %v", err)
	}

	cli, err := client.NewClient(client.WithTransport(httpTransport))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送JSON-RPC请求
	method := "getBlockByNumber"
	var params []interface{}
	result, err := cli.Request(ctx, method, params...)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	fmt.Printf("Result: %s\n", result)
}
