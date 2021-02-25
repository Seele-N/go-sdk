package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cmctype "github.com/cosmos/cosmos-sdk/codec/types"
)

// Codec amino codec to marshal/unmarshal
type Codec = codec.LegacyAmino

// AppCodec app marshaler
type AppCodec = codec.Marshaler

// NewCodec new codec
func NewCodec() (*Codec, AppCodec) {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := cmctype.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	return amino, marshaler
}
