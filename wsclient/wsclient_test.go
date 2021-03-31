package seelesdk

import (
	"encoding/hex"
	"fmt"
	"testing"

	//"github.com/Seele-N/go-sdk/wsclient"
	//"github.com/Seele-N/go-sdk/types"
	//"github.com/stretchr/testify/assert"
	"github.com/Seele-N/go-sdk/common"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	rpcURL = "tcp://127.0.0.1:26657"
)

func TestQueryStatus(t *testing.T) {
	common.SetConfig()
	valAddress := "seelevaloper1cdz5t52jhzgtp7z8et6xlly5nvkfm99fx7kmtm"
	addressBytes, _ := sdktypes.GetFromBech32(valAddress, common.ValidatorAddressPrefix)
	address, _ := sdktypes.AccAddressFromHex(hex.EncodeToString(addressBytes))
	fmt.Printf("address = %v", address.String())
	/*
		config, err := types.NewClientConfig(rpcURL, "seelehubv1", "broadcastMode", "0.01seele", 100000, 0, "")
		assert.Equal(t, nil, err)

		cli := NewWsClient(config)

		status, err := cli.Tendermint.QueryStatus()
		assert.Equal(t, nil, err)
		fmt.Printf("blockheight = %v", status.SyncInfo.LatestBlockHeight)
	*/
}
