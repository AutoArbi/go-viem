package wallet

import "github.com/AutoArbi/go-viem/client"

// Client is a Client for public Ethereum RPC methods
type Client struct {
	Client client.Transport
}
