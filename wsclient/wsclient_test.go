package seelesdk

import (
	"fmt"
	"testing"

	//"github.com/Seele-N/go-sdk/wsclient"
	"github.com/Seele-N/go-sdk/types"
	"github.com/stretchr/testify/assert"
)

const (
	rpcURL = "tcp://127.0.0.1:26657"
)

func TestQueryStatus(t *testing.T) {
	config, err := types.NewClientConfig(rpcURL, "seelehubv1", "broadcastMode", "0.01seele", 100000, 0, "")
	assert.Equal(t, nil, err)

	cli := NewWsClient(config)

	status, err := cli.Tendermint.QueryStatus()
	assert.Equal(t, nil, err)
	fmt.Printf("blockheight = %v", status.SyncInfo.LatestBlockHeight)
}
