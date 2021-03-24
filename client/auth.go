package client

import (
	"context"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// nolint
type (
	Account = authtypes.BaseAccount
)

// QueryAccount gets the account associated with an address on Seele
func (c *SeeleClient) QueryAccount(addr string) (acc Account, err error) {
	/*
		params := authtypes.QueryAccountRequest{
			Address: addr,
		}
		bz, err := c.cdc.MarshalJSON(params)
		if err != nil {
			return authtypes.BaseAccount{}, err
		}

		path := fmt.Sprintf("custom/%s/%s/%s", authtypes.QuerierRoute, authtypes.QueryAccount, addr)

		result, err := c.ABCIQuery(path, bz)
		if err != nil {
			return authtypes.BaseAccount{}, err
		}

		err = c.cdc.UnmarshalJSON(result, &acc)
		if err != nil {
			return authtypes.BaseAccount{}, err
		}

		return acc, err
	*/
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	authClient := authtypes.NewQueryClient(c.grpcConn)
	authRes, err := authClient.Account(
		ctx,
		&authtypes.QueryAccountRequest{Address: addr},
	)
	if err != nil {
		return authtypes.BaseAccount{}, err
	}

	err = acc.Unmarshal(authRes.Account.Value)
	if err != nil {
		return authtypes.BaseAccount{}, err
	}
	return acc, nil
}
