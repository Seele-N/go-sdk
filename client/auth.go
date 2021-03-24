package client

import (
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// nolint
type (
	Account = authtypes.BaseAccount
)

// QueryAccount gets the account associated with an address on Seele
func (c *SeeleClient) QueryAccount(addr string) (acc Account, err error) {
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
}
