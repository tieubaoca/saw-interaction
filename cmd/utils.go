package cmd

import (
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Seeds          []string `yaml:"seeds"`
	ChainId        string   `yaml:"chain_id"`
	AccountNumbers []uint64 `yaml:"account_numbers"`
	Sequences      []uint64 `yaml:"sequences"`
	Prefix         string   `yaml:"prefix"`
	Denom          string   `yaml:"denom"`
	From           string   `yaml:"from"`
	To             string   `yaml:"to"`
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

func readTxAndInitContexts(clientCtx client.Context, cmd *cobra.Command, filename string) (client.Context, tx.Factory, sdk.Tx, error) {
	log.Println("read tx from file")
	stdTx, err := authclient.ReadTxFromFile(clientCtx, filename)
	if err != nil {
		return clientCtx, tx.Factory{}, nil, err
	}
	log.Println("create tx factory")
	txFactory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		return clientCtx, tx.Factory{}, nil, err
	}
	log.Println("create output file")

	return clientCtx, txFactory, stdTx, nil
}

func signTx(cmd *cobra.Command, txCfg client.TxConfig, privKey types.PrivKey, cfg Config, txF tx.Factory, newTx sdk.Tx, filePath string) error {
	// f := cmd.Flags()
	// txCfg := clientCtx.TxConfig

	txBuilder, err := txCfg.WrapTxBuilder(newTx)

	if err != nil {
		return err
	}

	signature := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  txCfg.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: cfg.Sequences[0],
	}
	txBuilder.SetTimeoutHeight(999999999999)
	err = txBuilder.SetSignatures(signature)
	if err != nil {
		log.Fatal(err)
	}

	signature, err = tx.SignWithPrivKey(
		txCfg.SignModeHandler().DefaultMode(),
		authsigning.SignerData{
			ChainID:       cfg.ChainId,
			AccountNumber: cfg.AccountNumbers[0],
			Sequence:      cfg.Sequences[0],
		},
		txBuilder,
		privKey,
		txCfg,
		cfg.Sequences[0],
	)
	if err != nil {
		log.Fatal(err)
	}
	err = txBuilder.SetSignatures(signature)
	if err != nil {
		log.Fatal(err)
	}

	// set output
	// closeFunc, err := setOutputFile(cmd)
	if err != nil {
		return err
	}

	// defer closeFunc()
	// clientCtx.WithOutput(cmd.OutOrStdout())

	var json []byte

	json, err = marshalSignatureJSON(txCfg, txBuilder, false)
	if err != nil {
		return err
	}
	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	cmd.SetOut(fp)
	defer fp.Close()
	log.Println(string(json))
	cmd.Printf("%s\n", json)

	return err
}

func marshalSignatureJSON(txConfig client.TxConfig, txBldr client.TxBuilder, signatureOnly bool) ([]byte, error) {
	parsedTx := txBldr.GetTx()
	if signatureOnly {
		sigs, err := parsedTx.GetSignaturesV2()
		if err != nil {
			return nil, err
		}
		return txConfig.MarshalSignatureJSON(sigs)
	}

	return txConfig.TxJSONEncoder()(parsedTx)
}
