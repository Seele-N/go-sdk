package client

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// nolint
type (
	Coins = types.Coins
)

// QueryAccount gets the account associated with an address on Seele
func (c *SeeleClient) QueryAccountBank(addr string) (coins Coins, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	bankClient := banktypes.NewQueryClient(c.grpcConn)
	bankRes, err := bankClient.AllBalances(
		ctx,
		&banktypes.QueryAllBalancesRequest{Address: addr, Pagination: nil},
	)
	if err != nil {
		return Coins{}, err
	}

	return bankRes.GetBalances(), nil
}
