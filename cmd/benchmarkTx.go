/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// benchmarkTxCmd represents the benchmarkTx command
var benchmarkTxCmd = &cobra.Command{
	Use:   "benchmark-tx",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// numTx := cast.ToInt(cmd.Flag("NumTx").Value.String())
		grpcConn, err := grpc.Dial(cmd.Flag("host").Value.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		defer grpcConn.Close()
		// _ := txtypes.NewServiceClient(grpcConn)
		// var config Config
		// config.GetConfig()
		// interfaceRegistry := types.NewInterfaceRegistry()
		// marshaler := codec.NewProtoCodec(interfaceRegistry)
		// _ := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

		// for j, _ := range config.Seeds {
		// 	go func(j int) {
		// 		txChannel := make(chan []byte)
		// 		go TxChannelHandler(txClient, txChannel)
		// 		privKey := GetPrivKey(config.Seeds[j])

		// 		for i := 0; i < numTx; i++ {
		// 			txBuilder := NewBankSendTx(
		// 				txCfg,
		// 				privKey,
		// 				config.From,
		// 				config.To,
		// 				config.Denom,
		// 				config.ChainId,
		// 				config.Prefix,
		// 			)
		// 			txBytes, err := txCfg.TxEncoder()(txBuilder.GetTx())
		// 			if err != nil {
		// 				log.Fatal(err)
		// 			}
		// 			txChannel <- txBytes
		// 		}
		// 	}(j)

		// }
		// for {
		// }
	},
}

func init() {
	rootCmd.AddCommand(benchmarkTxCmd)

	// Here you will define your flags and configuration settings.
	benchmarkTxCmd.Flags().IntP("NumTx", "T", 500, "Number of txs want to send to blockchain node")
	benchmarkTxCmd.Flags().StringP("host", "H", "127.0.0.1:9090", "TM node to connect to")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// benchmarkTxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benchmarkTxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func TxChannelHandler(txClient txtypes.ServiceClient, txChannel chan []byte) {
	for {
		txBytes := <-txChannel
		BroadcastTx(txBytes, txClient)
	}
}

func BroadcastTx(txBytes []byte, txClient txtypes.ServiceClient) {
	log.Println("broadcast tx: ", hex.EncodeToString(txBytes))
	grpcRes, err := txClient.BroadcastTx(
		context.Background(),
		&txtypes.BroadcastTxRequest{
			Mode:    txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("broadcast tx response: ", grpcRes.TxResponse)
}
