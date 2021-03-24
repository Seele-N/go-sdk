package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/Seele-N/go-sdk/account"
	gosdktypes "github.com/Seele-N/go-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

const (
	maxABCIPathLength = 1024
	maxABCIDataLength = 1024 * 1024
)

var (
	ExceedABCIPathLengthError = fmt.Errorf("the abci path exceeds max length %d ", maxABCIPathLength)
	ExceedABCIDataLengthError = fmt.Errorf("the abci data exceeds max length %d ", maxABCIDataLength)
)

// SeeleClient seele client for Seele blockchain
type SeeleClient struct {
	rpc      rpcclient.Client
	wallet   account.Wallet
	config   *gosdktypes.ClientConfig
	cdc      *gosdktypes.Codec
	appCodec gosdktypes.AppCodec
}

// NewSeeleClient create a new SeeleClient
func NewSeeleClient(wallet account.Wallet, config gosdktypes.ClientConfig) *SeeleClient {
	rpc, err := rpchttp.New(config.NodeURI, "/websocket")
	if err != nil {
		panic(fmt.Sprintf("failed to get client: %s", err))
	}
	cdc, appCodec := gosdktypes.NewCodec()
	client := &SeeleClient{
		rpc:      rpc,
		wallet:   wallet,
		config:   &config,
		cdc:      cdc,
		appCodec: appCodec,
	}

	return client
}

// // ABCIQuery sends a query to Seele
func (sc *SeeleClient) ABCIQuery(path string, data tmbytes.HexBytes) ([]byte, error) {
	if err := ValidateABCIQuery(path, data); err != nil {
		return []byte{}, err
	}

	result, err := sc.rpc.ABCIQuery(context.Background(), path, data)
	if err != nil {
		return []byte{}, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return []byte{}, errors.New(resp.Log)
	}

	value := result.Response.GetValue()
	if len(value) == 0 {
		return []byte{}, nil
	}

	return value, nil
}

// ValidateABCIQuery validates an ABCI query
func ValidateABCIQuery(path string, data tmbytes.HexBytes) error {
	if err := ValidateABCIPath(path); err != nil {
		return err
	}
	if err := ValidateABCIData(data); err != nil {
		return err
	}
	return nil
}

// ValidateABCIPath validates an ABCI query's path
func ValidateABCIPath(path string) error {
	if len(path) > maxABCIPathLength {
		return ExceedABCIPathLengthError
	}
	return nil
}

// ValidateABCIData validates an ABCI query's data
func ValidateABCIData(data tmbytes.HexBytes) error {
	if len(data) > maxABCIDataLength {
		return ExceedABCIPathLengthError
	}
	return nil
}
