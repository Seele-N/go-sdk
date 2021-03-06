package seelesdk

import (
	"fmt"

	gosdktypes "github.com/Seele-N/go-sdk/types"
	"github.com/Seele-N/go-sdk/wsclient/tendermint"

	//tmbytes "github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

// WsClient websocket client
type WsClient struct {
	rpc      rpcclient.Client
	config   *gosdktypes.ClientConfig
	cdc      *gosdktypes.Codec
	appCodec gosdktypes.AppCodec

	Tendermint tendermint.ClientService
}

// NewWsClient new websocket client
func NewWsClient(config gosdktypes.ClientConfig) *WsClient {
	rpc, err := rpchttp.New(config.NodeURI, "/websocket")
	if err != nil {
		panic(fmt.Sprintf("failed to get client: %s", err))
	}
	cdc, appCodec := gosdktypes.NewCodec()
	client := &WsClient{
		rpc:      rpc,
		config:   &config,
		cdc:      cdc,
		appCodec: appCodec,
	}

	client.Tendermint = tendermint.NewClient(rpc)

	return client
}

/*
// Query executes the basic query
func (wc *WsClient) Query(path string, key tmbytes.HexBytes) (res []byte, height int64, err error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	}

	result, err := wc.ABCIQueryWithOptions(context.Background(), path, key, opts)
	if err != nil {
		return
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, height, errors.New(resp.Log)
	}

	return resp.Value, resp.Height, err
}
*/
