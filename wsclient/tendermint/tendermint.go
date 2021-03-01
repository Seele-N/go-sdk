package tendermint

import (
	"context"
	"encoding/hex"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// nolint
type (
	Block              = tmtypes.Block
	ResultBlock        = ctypes.ResultBlock
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
	QueryStatus() (pStatus *ResultStatus, err error)
	QueryBlock(height int64) (*ResultBlock, error)
	QueryBlockResults(height int64) (*ResultBlockResults, error)
	QueryBlockByHash(hashHexStr string) (*ResultBlock, error)
	QueryTxResult(hashHexStr string, prove bool) (pResultTx *ResultTx, err error)
	QueryValidatorsResult(height int64) (pValsResult *ResultValidators, err error)
	QueryCommitResult(height int64) (pCommitResult *ResultCommit, err error)
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
func (c Client) QueryBlock(height int64) (pBlock *ResultBlock, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	return c.Block(context.Background(), pHeight)
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

// QueryBlockByHash get the abci result of the block by hash
func (c Client) QueryBlockByHash(hashHexStr string) (pBlock *ResultBlock, err error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return
	}

	return c.BlockByHash(context.Background(), hash)
}

// QueryTxResult gets the detail info of a tx with its tx hash
func (c Client) QueryTxResult(hashHexStr string, prove bool) (pResultTx *ResultTx, err error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return
	}

	return c.Tx(context.Background(), hash, prove)
}

// QueryValidatorsResult gets the validators info on a specific height
// query the latest block with height 0 input
func (c Client) QueryValidatorsResult(height int64) (pValsResult *ResultValidators, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	var page, perPage int
	page = 1
	perPage = 0
	return c.Validators(context.Background(), pHeight, &page, &perPage)
}

// QueryCommitResult gets the commit info of the block on a specific height
// query the latest block with height 0 input
func (c Client) QueryCommitResult(height int64) (pCommitResult *ResultCommit, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	return c.Commit(context.Background(), pHeight)
}
