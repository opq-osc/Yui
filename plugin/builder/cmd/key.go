/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/charmbracelet/log"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "生成公，私钥",
	Long:  `生成公，私钥`,
	RunE: func(cmd *cobra.Command, args []string) error {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			return err
		}
		privateKeyBytes := crypto.FromECDSA(privateKey)
		log.Info("私钥", "key", hex.EncodeToString(privateKeyBytes))
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
		log.Info("公钥", "publicKey", hex.EncodeToString(publicKeyBytes))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(keyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
