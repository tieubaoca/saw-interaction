/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		httpConn, err := connect(cmd.Flag("host").Value.String())
		if err != nil {
			return err
		}
		ctxTO, _ := context.WithTimeout(context.Background(), 10*time.Second)
		bcInfo, err := httpConn.BlockchainInfo(ctxTO, 0, 1)
		if err != nil {
			return err
		}
		fmt.Printf("connected successfully. last height: %v\n", bcInfo.LastHeight)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.
	connectCmd.Flags().StringP("host", "H", "tcp://localhost:26657", "TM node to connect to")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func connect(host string) (*http.HTTP, error) {
	httpConn, err := client.NewClientFromNode(host)
	if err != nil {
		return nil, err
	}
	return httpConn, nil
}
