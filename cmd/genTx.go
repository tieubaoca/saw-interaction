/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// genTxCmd represents the genTx command
var genTxCmd = &cobra.Command{
	Use:   "gen-tx",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config Config
		config.GetConfig()
		privKey := GetPrivKey(config.Seeds[0])
		bench32Addr, err := bech32.ConvertAndEncode(config.Prefix, privKey.PubKey().Address())
		if err != nil {
			panic(err)
		}
		cmd.Println("Address: " + bench32Addr)

		res, err := http.Get("https://docs-api.onechain.game/cosmos/auth/v1beta1/account_info/" + bench32Addr)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		responseData, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		cmd.Println("Response: " + string(responseData))
		var accInfo map[string]interface{}
		err = json.Unmarshal(responseData, &accInfo)
		if err != nil {
			panic(err)
		}
		accSeq := cast.ToUint64(accInfo["info"].(map[string]interface{})["sequence"].(string))
		accNum := cast.ToUint64(accInfo["info"].(map[string]interface{})["account_number"].(string))

		interfaceRegistry := types.NewInterfaceRegistry()
		marshaler := codec.NewProtoCodec(interfaceRegistry)
		txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)
		txBuilder := NewBankSendTx(
			txCfg,
			privKey,
			accNum,
			accSeq,
			config.Denom,
			config.ChainId,
			config.Prefix,
		)
		txBytes, err := txCfg.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			panic(err)
		}
		cmd.Println("Tx bytes: 0x" + hex.EncodeToString(txBytes))

		txres, err := http.Get("https://api-chain.onechain.game/broadcast_tx_commit?tx=0x" + hex.EncodeToString(txBytes))
		if err != nil {
			panic(err)
		}
		defer txres.Body.Close()
		txResponseData, err := io.ReadAll(txres.Body)
		if err != nil {
			panic(err)
		}
		cmd.Println("Tx Response: " + string(txResponseData))
	},
}

func init() {
	rootCmd.AddCommand(genTxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genTxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genTxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
