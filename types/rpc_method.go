package types

type ETHMethod string

// Block api
const (
	GetBlockByNumber                 ETHMethod = "eth_getBlockByNumber"
	GetBlockByHash                   ETHMethod = "eth_getBlockByHash"
	GetBlockNumber                   ETHMethod = "eth_blockNumber"
	GetBlockTransactionCountByHash   ETHMethod = "eth_getBlockTransactionCountByHash"
	GetBlockTransactionCountByNumber ETHMethod = "eth_getBlockTransactionCountByNumber"
	FeeHistory                       ETHMethod = "eth_feeHistory"
	GetUncleByBlockHashAndIndex      ETHMethod = "eth_getUncleByBlockHashAndIndex"
	GetUncleByBlockNumberAndIndex    ETHMethod = "eth_getUncleByBlockNumberAndIndex"
	GetUncleCountByBlockHash         ETHMethod = "eth_getUncleCountByBlockHash"
	GetUncleCountByBlockNumber       ETHMethod = "eth_getUncleCountByBlockNumber"
)

// Transaction  API
const (
	GetTransactionByHash                ETHMethod = "eth_getTransactionByHash"
	GetTransactionByBlockHashAndIndex   ETHMethod = "eth_getTransactionByBlockHashAndIndex"
	GetTransactionByBlockNumberAndIndex ETHMethod = "eth_getTransactionByBlockNumberAndIndex"
	GetTransactionReceipt               ETHMethod = "eth_getTransactionReceipt"
	SendRawTransaction                  ETHMethod = "eth_sendRawTransaction"
	PendingTransactions                 ETHMethod = "eth_pendingTransactions"
	GetBlockReceipts                    ETHMethod = "eth_getBlockReceipts"
)

// account / contract API
const (
	GetBalance   ETHMethod = "eth_getBalance"
	GetCode      ETHMethod = "eth_getCode"
	GetStorageAt ETHMethod = "eth_getStorageAt"
	Call         ETHMethod = "eth_call"
	EstimateGas  ETHMethod = "eth_estimateGas"
	Accounts     ETHMethod = "eth_accounts"
	Coinbase     ETHMethod = "eth_coinbase"
)

// 链状态 & 其他 API
const (
	ChainID         ETHMethod = "eth_chainId"
	GasPrice        ETHMethod = "eth_gasPrice"
	Syncing         ETHMethod = "eth_syncing"
	ProtocolVersion ETHMethod = "eth_protocolVersion"
)

// NetworkMethod net api
type NetworkMethod string

const (
	NetVersion   NetworkMethod = "net_version"
	NetListening NetworkMethod = "net_listening"
	NetPeerCount NetworkMethod = "net_peerCount"
)

// Web3Method Web3 api
type Web3Method string

const (
	Web3ClientVersion Web3Method = "web3_clientVersion"
	Web3Sha3          Web3Method = "web3_sha3"
)
