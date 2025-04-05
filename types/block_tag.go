package types

type BlockTag string

const (
	// EARLIEST Genesis Block
	EARLIEST BlockTag = "earliest"

	// LATEST proposed blocks
	LATEST BlockTag = "latest"

	// SAFE The latest and safest header block
	SAFE BlockTag = "safe"

	// FINALIZED Latest finalized block
	FINALIZED BlockTag = "finalized"

	// PENDING status/Transaction
	PENDING BlockTag = "pending"
)
