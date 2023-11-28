/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/cobra"
)

// genTxCmd represents the genTx command
var genTxCmd = &cobra.Command{
	Use:   "gen-tx [file-path]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		var config Config
		config.GetConfig()
		privKey := GetPrivKey(config.Seeds[0])
		bench32Addr, err := bech32.ConvertAndEncode(config.Prefix, privKey.PubKey().Address())
		if err != nil {
			panic(err)
		}
		cmd.Println("Address: " + bench32Addr)

		// res, err := http.Get("http://localhost:1317/cosmos/auth/v1beta1/account_info/" + bench32Addr)
		// if err != nil {
		// 	panic(err)
		// }
		// defer res.Body.Close()
		// responseData, err := io.ReadAll(res.Body)
		// if err != nil {
		// 	panic(err)
		// }
		// cmd.Println("Response: " + string(responseData))
		// var accInfo map[string]interface{}
		// err = json.Unmarshal(responseData, &accInfo)
		// if err != nil {
		// 	panic(err)
		// }
		// accSeq := cast.ToUint64(accInfo["info"].(map[string]interface{})["sequence"].(string))
		// accNum := cast.ToUint64(accInfo["info"].(map[string]interface{})["account_number"].(string))
		amino := codec.NewLegacyAmino()
		interfaceRegistry := types.NewInterfaceRegistry()
		marshaler := codec.NewProtoCodec(interfaceRegistry)
		txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

		std.RegisterLegacyAminoCodec(amino)
		std.RegisterInterfaces(interfaceRegistry)
		ModuleBasics.RegisterLegacyAminoCodec(amino)
		ModuleBasics.RegisterInterfaces(interfaceRegistry)
		// txBuilder := NewBankSendTx(
		// 	txCfg,
		// 	privKey,
		// 	config.From,
		// 	config.To,
		// 	config.AccountNumbers[0],
		// 	config.Sequences[0],
		// 	config.Denom,
		// 	config.ChainId,
		// 	config.Prefix,
		// )
		_, txF, tx, err := readTxAndInitContexts(client.Context{
			FromAddress:       sdk.AccAddress(bench32Addr),
			From:              bench32Addr,
			ChainID:           config.ChainId,
			InterfaceRegistry: interfaceRegistry,
			Codec:             marshaler,
			TxConfig:          txCfg,
			LegacyAmino:       amino,
		}, cmd, filePath)
		if err != nil {
			log.Fatal(err)
		}

		err = signTx(cmd, txCfg, privKey, config, txF, tx, filePath)
		if err != nil {
			log.Fatal(err)
		}

		// cmd.Println("Tx bytes: 0x" + hex.EncodeToString(txBytes))

		// txres, err := http.Get("https://api-chain.onechain.game/broadcast_tx_commit?tx=0x" + hex.EncodeToString(txBytes))
		// if err != nil {
		// 	panic(err)
		// }
		// defer txres.Body.Close()
		// txResponseData, err := io.ReadAll(txres.Body)
		// if err != nil {
		// 	panic(err)
		// }
		// cmd.Println("Tx Response: " + string(txResponseData))
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
