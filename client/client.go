package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Seele-N/go-sdk/account"
	gosdktypes "github.com/Seele-N/go-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	signing2 "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	"google.golang.org/grpc"
)

const (
	maxABCIPathLength = 1024
	maxABCIDataLength = 1024 * 1024
)

var (
	ExceedABCIPathLengthError = fmt.Errorf("the abci path exceeds max length %d ", maxABCIPathLength)
	ExceedABCIDataLengthError = fmt.Errorf("the abci data exceeds max length %d ", maxABCIDataLength)
)

// SeeleClient seele client for Seele blockchain
type SeeleClient struct {
	rpc      rpcclient.Client
	grpcConn *grpc.ClientConn
	wallet   account.Wallet
	config   *gosdktypes.ClientConfig
	cdc      *gosdktypes.Codec
	appCodec gosdktypes.AppCodec
	txCfg    client.TxConfig
}

// NewSeeleClient create a new SeeleClient
func NewSeeleClient(wallet account.Wallet, config gosdktypes.ClientConfig) *SeeleClient {

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		config.GRpcURI, // your gRPC server address.
		grpc.WithInsecure(),
		grpc.WithBlock(), // The SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		panic(fmt.Sprintf("failed to get grpc client: %s", err))
	}

	rpc, err := rpchttp.New(config.NodeURI, "/websocket")
	if err != nil {
		panic(fmt.Sprintf("failed to get rpc client: %s", err))
	}
	cdc, appCodec := gosdktypes.NewCodec()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

	client := &SeeleClient{
		rpc:      rpc,
		grpcConn: grpcConn,
		wallet:   wallet,
		config:   &config,
		cdc:      cdc,
		appCodec: appCodec,
		txCfg:    txCfg,
	}

	return client
}

// Wallet return client wallet account
func (sc *SeeleClient) Wallet() account.Wallet {
	return sc.wallet
}

// ABCIQuery sends a query to Seele
func (sc *SeeleClient) ABCIQuery(path string, data tmbytes.HexBytes) ([]byte, error) {
	if err := ValidateABCIQuery(path, data); err != nil {
		return []byte{}, err
	}

	result, err := sc.rpc.ABCIQuery(context.Background(), path, data)
	if err != nil {
		return []byte{}, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return []byte{}, errors.New(resp.Log)
	}

	value := result.Response.GetValue()
	if len(value) == 0 {
		return []byte{}, nil
	}

	return value, nil
}

// ValidateABCIQuery validates an ABCI query
func ValidateABCIQuery(path string, data tmbytes.HexBytes) error {
	if err := ValidateABCIPath(path); err != nil {
		return err
	}
	if err := ValidateABCIData(data); err != nil {
		return err
	}
	return nil
}

// ValidateABCIPath validates an ABCI query's path
func ValidateABCIPath(path string) error {
	if len(path) > maxABCIPathLength {
		return ExceedABCIPathLengthError
	}
	return nil
}

// ValidateABCIData validates an ABCI query's data
func ValidateABCIData(data tmbytes.HexBytes) error {
	if len(data) > maxABCIDataLength {
		return ExceedABCIPathLengthError
	}
	return nil
}

func (c *SeeleClient) buildandsend(msgs []sdk.Msg) {

	tx := c.txCfg.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}

	fee := sdk.NewCoins(sdk.NewCoin("snp", sdk.NewInt(0)))

	tx.SetMemo("")
	tx.SetFeeAmount(fee)
	tx.SetGasLimit(1000000)

	account, _ := c.QueryAccount(sdk.AccAddress(c.Wallet().GetPubAddress()).String())
	sequence := account.Sequence
	number := account.AccountNumber

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	priv := c.Wallet().GetPrivKey()
	sig := signing.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &sigData,
		Sequence: sequence,
	}
	if err := tx.SetSignatures(sig); err != nil {
		panic(err)
	}

	signBytes, err := c.txCfg.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signing2.SignerData{
		ChainID:       c.config.ChainID,
		AccountNumber: number,
		Sequence:      sequence,
	}, tx.GetTx())
	if err != nil {
		panic(err)
	}

	// Sign those bytes
	sigBytes, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &sigData,
		Sequence: sequence,
	}

	if err := tx.SetSignatures(sig); err != nil {
		panic(err)
	}

	txBytes, err := c.txCfg.TxEncoder()(tx.GetTx())
	if err != nil {
		panic(err)
	}

	c.sendTx(txBytes)
}

func (sc *SeeleClient) sendTx(txBytes []byte) {
	txClient := tx.NewServiceClient(sc.grpcConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	txClient.BroadcastTx(
		ctx,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_ASYNC,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
		//grpc.WaitForReady(false), grpc.FailFast(false),
	)
}
