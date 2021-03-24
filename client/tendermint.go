import client

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

// QueryChainID query chain id
func (c *SeeleClient) QueryChainID() (string,error){
	result,err := c.QueryStatus()
	if err != nil{
		return "",err
	}
	return result.NodeInfo.Network, nil
}

// QueryStatus gets the blockchain info
func (c *SeeleClient) QueryStatus() (pStatus *ResultStatus, err error) {

	return c.rpc.Status(context.Background())
}

// QueryBlock gets the block info of a specific height
// query the latest block with height 0 input
func (c *SeeleClient) QueryBlock(height int64) (pBlock *ResultBlock, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	return c.rpc.Block(context.Background(), pHeight)
}

// QueryBlockResults gets the abci result of the block on a specific height
// query the latest block with height 0 input
func (c *SeeleClient) QueryBlockResults(height int64) (pBlockResults *ResultBlockResults, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	return c.rpc.BlockResults(context.Background(), pHeight)
}

// QueryBlockByHash get the abci result of the block by hash
func (c *SeeleClient) QueryBlockByHash(hashHexStr string) (pBlock *ResultBlock, err error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return
	}

	return c.rpc.BlockByHash(context.Background(), hash)
}

// QueryTxResult gets the detail info of a tx with its tx hash
func (c *SeeleClient) QueryTxResult(hashHexStr string, prove bool) (pResultTx *ResultTx, err error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return
	}

	return c.rpc.Tx(context.Background(), hash, prove)
}

// QueryValidatorsResult gets the validators info on a specific height
// query the latest block with height 0 input
func (c *SeeleClient) QueryValidatorsResult(height int64) (pValsResult *ResultValidators, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	var page, perPage int
	page = 1
	perPage = 0
	return c.rpc.Validators(context.Background(), pHeight, &page, &perPage)
}

// QueryCommitResult gets the commit info of the block on a specific height
// query the latest block with height 0 input
func (c *SeeleClient) QueryCommitResult(height int64) (pCommitResult *ResultCommit, err error) {

	var pHeight *int64
	if height > 0 {
		pHeight = &height
	}

	return c.rpc.Commit(context.Background(), pHeight)
}