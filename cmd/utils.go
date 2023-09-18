package cmd

import (
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Seeds          []string `yaml:"seeds"`
	AccountNumbers []uint64 `yaml:"account_numbers"`
	Sequences      []uint64 `yaml:"sequences"`
	ChainId        string   `yaml:"chain_id"`
	Prefix         string   `yaml:"prefix"`
	Denom          string   `yaml:"denom"`
}

func (c *Config) GetConfig() *Config {
	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func GetPrivKey(seed string) types.PrivKey {
	hdPath := hd.CreateHDPath(sdk.GetConfig().GetCoinType(), 0, 0).String()
	bip39Passphrase := ""

	derivedPriv, err := hd.Secp256k1.Derive()(seed, bip39Passphrase, hdPath)
	if err != nil {
		log.Fatal(err)
	}
	privKey := hd.Secp256k1.Generate()(derivedPriv)
	return privKey
}

func NewBankSendTx(txCfg client.TxConfig, privKey types.PrivKey, accountNumber, sequence uint64, denom, chainId, prefix string) client.TxBuilder {
	txBuilder := txCfg.NewTxBuilder()
	sender, err := bech32.ConvertAndEncode(prefix, privKey.PubKey().Address())
	if err != nil {
		log.Fatal(err)
	}
	txBuilder.SetMsgs(
		&banktypes.MsgSend{
			FromAddress: sender,
			ToAddress:   sender,
			Amount:      sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1))),
		},
	)

	signature := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  txCfg.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: sequence,
	}
	// txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(20))))
	txBuilder.SetGasLimit(200000)
	txBuilder.SetMemo("")
	txBuilder.SetTimeoutHeight(999999999999)
	err = txBuilder.SetSignatures(signature)
	if err != nil {
		log.Fatal(err)
	}

	signature, err = tx.SignWithPrivKey(
		txCfg.SignModeHandler().DefaultMode(),
		authsigning.SignerData{
			ChainID:       chainId,
			AccountNumber: accountNumber,
			Sequence:      sequence,
		},
		txBuilder,
		privKey,
		txCfg,
		sequence,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = txBuilder.SetSignatures(signature)
	if err != nil {
		log.Fatal(err)
	}
	return txBuilder
}
