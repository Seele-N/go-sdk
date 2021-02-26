package tendermint

import (
	"context"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// nolint
type (
	Block              = tmtypes.Block
	ResultBlockResults = ctypes.ResultBlockResults
	ResultCommit       = ctypes.ResultCommit
	ResultValidators   = ctypes.ResultValidators
	ResultTx           = ctypes.ResultTx
	ResultTxSearch     = ctypes.ResultTxSearch
	ResultStatus       = ctypes.ResultStatus
)

// Client client for tendermint
type Client struct {
	rpcclient.Client
}

// ClientService the interface of tendermint client
type ClientService interface {
	QueryBlock(height int64) (*Block, error)
	QueryBlockResults(height int64) (*ResultBlockResults, error)
	QueryStatus() (pStatus *ResultStatus, err error)
}

// NewClient new tendermint client
func NewClient(rpc rpcclient.Client) ClientService {
	return Client{
		Client: rpc,
	}
}

// QueryStatus gets the blockchain info
func (c Client) QueryStatus() (pStatus *ResultStatus, err error) {

	return c.Status(context.Background())
}

// QueryBlock gets the block info of a specific height
// query the latest block with height 0 input
func (c Client) QueryBlock(height int64) (pBlock *Block, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	pTmBlockResult, err := c.Block(context.Background(), pHeight)
	if err != nil {
		return
	}

	return pTmBlockResult.Block, err
}

// QueryBlockResults gets the abci result of the block on a specific height
// query the latest block with height 0 input
func (c Client) QueryBlockResults(height int64) (pBlockResults *ResultBlockResults, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	return c.BlockResults(context.Background(), pHeight)
}
