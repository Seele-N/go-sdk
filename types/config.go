package types

//sdk "github.com/cosmos/cosmos-sdk/types"

// ClientConfig records the base config of gosdk client
type ClientConfig struct {
	NodeURI       string
	ChainID       string
	BroadcastMode string
	Gas           uint64
	GasAdjustment float64
	//Fees          sdk.DecCoins
	//GasPrices     sdk.DecCoins
}

// NewClientConfig new ClientConfig
func NewClientConfig(nodeURI, chainID string, broadcastMode string, feesStr string, gas uint64, gasAdjustment float64,
	gasPricesStr string) (cliConfig ClientConfig, err error) {

	return ClientConfig{
		NodeURI:       nodeURI,
		ChainID:       chainID,
		BroadcastMode: broadcastMode,
		Gas:           gas,
		GasAdjustment: gasAdjustment,
		//Fees:          fees,
		//GasPrices:     gasPrices,
	}, err
}
