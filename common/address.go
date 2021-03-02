package common

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
)

// GetPubKeyFromString get pubkey from string
func GetPubKeyFromString(pkstr string) (cryptotypes.PubKey, error) {
	bz, err := hex.DecodeString(pkstr)
	if err == nil {
		if len(bz) == ed25519.PubKeySize {
			return &ed25519.PubKey{Key: bz}, nil
		}
	}

	bz, err = base64.StdEncoding.DecodeString(pkstr)
	if err == nil {
		if len(bz) == ed25519.PubKeySize {
			return &ed25519.PubKey{Key: bz}, nil
		}
	}

	pk, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, pkstr)
	if err == nil {
		return pk, nil
	}

	pk, err = types.GetPubKeyFromBech32(types.Bech32PubKeyTypeValPub, pkstr)
	if err == nil {
		return pk, nil
	}

	pk, err = types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, pkstr)
	if err == nil {
		return pk, nil
	}

	return nil, fmt.Errorf("pubkey '%s' invalid; expected hex, base64, or bech32 of correct size", pkstr)
}

// GetAccountPubKey Get Account PubKey from ed25519 public key
func GetAccountPubKey(pk cryptotypes.PubKey) (string, error) {
	return types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pk)
}

// GetValidatorPubKey Get Validator PubKey from ed25519 public key
func GetValidatorPubKey(pk cryptotypes.PubKey) (string, error) {
	return types.Bech32ifyPubKey(types.Bech32PubKeyTypeValPub, pk)
}

// GetConsenusPubKey Get Consenus PubKey from ed25519 public key
func GetConsenusPubKey(pk cryptotypes.PubKey) (string, error) {
	return types.Bech32ifyPubKey(types.Bech32PubKeyTypeConsPub, pk)
}
