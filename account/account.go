package account

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Seele-N/go-sdk/common"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"

	//ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Wallet interface {
	GetPrivKey() crypto.PrivKey
	GetPubAddress() crypto.Address
	ExportAsMnemonic() (string, error)
	ExportAsPrivateKey() (string, error)
	//ExportAsKeyStore(password string) (*EncryptedKeyJSON, error)
}

type wallet struct {
	privKey crypto.PrivKey
	//addr     ctypes.AccAddress
	mnemonic string
}

func (w *wallet) GetPrivKey() crypto.PrivKey {
	return w.privKey
}

func (w *wallet) GetPubAddress() crypto.Address {
	return w.privKey.PubKey().Address()
}

func (w *wallet) ExportAsMnemonic() (string, error) {
	if w.mnemonic == "" {
		return "", fmt.Errorf("This wallet is not recover from mnemonic or anto generated ")
	}
	return w.mnemonic, nil
}

func (w *wallet) ExportAsPrivateKey() (string, error) {
	secpPrivateKey, ok := w.privKey.(secp256k1.PrivKey)
	if !ok {
		return "", fmt.Errorf(" Only PrivKeySecp256k1 key is supported ")
	}
	return hex.EncodeToString(secpPrivateKey[:]), nil
}

/*
func (w *wallet) ExportAsKeyStore(password string) (*EncryptedKeyJSON, error) {
	return generateKeyStore(w.GetPrivKey(), password)
}
*/

func NewWallet() (Wallet, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	return NewWalletFromMnemonic(mnemonic)
}

func NewWalletFromMnemonic(mnemonic string) (Wallet, error) {
	w := wallet{}
	err := w.createFromMnemonic(mnemonic, common.FullFundraiserPath)
	return &w, err
}

func NewWalletFromPrivateKey(priKey string) (Wallet, error) {
	w := wallet{}
	err := w.createFromPrivateKey(priKey)
	return &w, err
}

func NewWalletFromKeyStore(file string, password string) (Wallet, error) {
	w := wallet{}
	err := w.createFromKeyStore(file, password)
	return &w, err
}

func (w *wallet) createFromMnemonic(mnemonic, keyPath string) error {
	words := strings.Split(mnemonic, " ")
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("mnemonic length should either be 12 or 24")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return err
	}

	// create master key and derive first key:
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, keyPath)
	if err != nil {
		return err
	}
	priKey := secp256k1.PrivKey(derivedPriv)
	//addr := priKey.PubKey().Address()
	//m.addr = addr
	w.privKey = priKey
	w.mnemonic = mnemonic

	return nil
}

func (w *wallet) createFromPrivateKey(privateKey string) error {
	priBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}

	if len(priBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], priBytes[:32])
	priKey := secp256k1.PrivKey(keyBytesArray[:])
	//addr := ctypes.AccAddress(priKey.PubKey().Address())
	//w.addr = addr
	w.privKey = priKey
	return nil
}

func (w *wallet) createFromKeyStore(keystoreFile string, password string) error {
	if password == "" {
		return fmt.Errorf("Password is missing ")
	}
	keyJson, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return err
	}
	var encryptedKey EncryptedKeyJSON
	err = json.Unmarshal(keyJson, &encryptedKey)
	if err != nil {
		return err
	}
	keyBytes, err := decryptKey(&encryptedKey, password)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	priKey := secp256k1.PrivKey(keyBytesArray[:])
	//addr := ctypes.AccAddress(priKey.PubKey().Address())
	//m.addr = addr
	w.privKey = priKey
	return nil
}
