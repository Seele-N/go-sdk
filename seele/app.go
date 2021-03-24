package seele

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/types"
)

func MakeCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()

	codec.RegisterEvidences(cdc)

	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterInterface((*types.Tx)(nil), nil)
	cryptocodec.RegisterCrypto(cdc)

	return cdc
}
