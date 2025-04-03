package public

import "github.com/AutoArbi/go-viem/client"

// Client is a client for public Ethereum RPC methods
type Client struct {
	client client.Interface
}
