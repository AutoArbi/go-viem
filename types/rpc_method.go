package types

type RPCMethod string

// Block api
const (
	GetBlockByNumber                 RPCMethod = "eth_getBlockByNumber"
	GetBlockByHash                   RPCMethod = "eth_getBlockByHash"
	GetBlockNumber                   RPCMethod = "eth_blockNumber"
	GetBlockTransactionCountByHash   RPCMethod = "eth_getBlockTransactionCountByHash"
	GetBlockTransactionCountByNumber RPCMethod = "eth_getBlockTransactionCountByNumber"
	FeeHistory                       RPCMethod = "eth_feeHistory"
	GetUncleByBlockHashAndIndex      RPCMethod = "eth_getUncleByBlockHashAndIndex"
	GetUncleByBlockNumberAndIndex    RPCMethod = "eth_getUncleByBlockNumberAndIndex"
	GetUncleCountByBlockHash         RPCMethod = "eth_getUncleCountByBlockHash"
	GetUncleCountByBlockNumber       RPCMethod = "eth_getUncleCountByBlockNumber"

	WatchBlocks      RPCMethod = "eth_watchBlocks"
	WatchBlockNumber RPCMethod = "eth_watchBlockNumber"

	SimulateCall   RPCMethod = "eth_simulateCall"
	SimulateBlocks RPCMethod = "eth_simulateBlocks"
)

// Transaction  API
const (
	GetTransactionCount                 RPCMethod = "eth_getTransactionCount"
	GetTransactionByHash                RPCMethod = "eth_getTransactionByHash"
	GetTransactionByBlockHashAndIndex   RPCMethod = "eth_getTransactionByBlockHashAndIndex"
	GetTransactionByBlockNumberAndIndex RPCMethod = "eth_getTransactionByBlockNumberAndIndex"
	GetTransactionReceipt               RPCMethod = "eth_getTransactionReceipt"
	SendRawTransaction                  RPCMethod = "eth_sendRawTransaction"
	PendingTransactions                 RPCMethod = "eth_pendingTransactions"
	GetBlockReceipts                    RPCMethod = "eth_getBlockReceipts"
)

// account / contract API
const (
	GetBalance   RPCMethod = "eth_getBalance"
	GetChainID   RPCMethod = "eth_chainId"
	GetCode      RPCMethod = "eth_getCode"
	GetStorageAt RPCMethod = "eth_getStorageAt"
	Call         RPCMethod = "eth_call"
	EstimateGas  RPCMethod = "eth_estimateGas"
	Accounts     RPCMethod = "eth_accounts"
	Coinbase     RPCMethod = "eth_coinbase"
)

// chain API / other API
const (
	CreateAccessList RPCMethod = "eth_createAccessList"
	GasPrice         RPCMethod = "eth_gasPrice"
	Syncing          RPCMethod = "eth_syncing"
	ProtocolVersion  RPCMethod = "eth_protocolVersion"
)

// net api
const (
	NetVersion   RPCMethod = "net_version"
	NetListening RPCMethod = "net_listening"
	NetPeerCount RPCMethod = "net_peerCount"
)

// Web3 api
const (
	Web3ClientVersion RPCMethod = "web3_clientVersion"
	Web3Sha3          RPCMethod = "web3_sha3"
)

// logs / event API
const (
	GetLogs                RPCMethod = "eth_getLogs"
	GetLogByHash           RPCMethod = "eth_getLogByHash"
	TransactionSubscribe   RPCMethod = "eth_subscribe"
	TransactionUnsubscribe RPCMethod = "eth_unsubscribe"
)
