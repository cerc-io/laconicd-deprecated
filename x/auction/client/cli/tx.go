package cli

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tharsis/ethermint/x/auction/types"

	wnsUtils "github.com/tharsis/ethermint/utils"
)

// GetTxCmd returns transaction commands for this module.
func GetTxCmd() *cobra.Command {
	auctionTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auction transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	auctionTxCmd.AddCommand(
		GetCmdCreateAuction(),
		GetCmdCommitBid(),
		GetCmdRevealBid(),
	)

	return auctionTxCmd
}

func GetCmdCreateAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [commits-duration] [reveals-duration] [commit-fee] [reveal-fee] [minimum-bid]",
		Short: "Create auction.",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			commitsDuration, err := time.ParseDuration(args[0])
			if err != nil {
				return err
			}

			revealsDuration, err := time.ParseDuration(args[1])
			if err != nil {
				return err
			}

			commitFee, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			revealFee, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			minimumBid, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return err
			}

			params := types.Params{
				CommitsDuration: commitsDuration,
				RevealsDuration: revealsDuration,
				CommitFee:       commitFee,
				RevealFee:       revealFee,
				MinimumBid:      minimumBid,
			}
			msg := types.NewMsgCreateAuction(params, clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCommitBid is the CLI command for committing a bid.
func GetCmdCommitBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-bid [auction-id] [bid-amount]",
		Short: "Commit sealed bid.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bidAmount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			mnemonic, err := wnsUtils.GenerateMnemonic()
			if err != nil {
				return err
			}

			chainID := viper.GetString("chain-id")
			auctionID := args[0]

			reveal := map[string]interface{}{
				"chainId":       chainID,
				"auctionId":     auctionID,
				"bidderAddress": clientCtx.GetFromAddress().String(),
				"bidAmount":     bidAmount.String(),
				"noise":         mnemonic,
			}

			commitHash, content, err := wnsUtils.GenerateHash(reveal)
			if err != nil {
				return err
			}

			// Save reveal file.
			ioutil.WriteFile(fmt.Sprintf("%s-%s.json", clientCtx.GetFromName(), commitHash), content, 0o600)

			msg := types.NewMsgCommitBid(auctionID, commitHash, clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRevealBid is the CLI command for revealing a bid.
func GetCmdRevealBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reveal-bid [auction-id] [reveal-file-path]",
		Short: "Reveal bid.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionID := args[0]
			revealFilePath := args[1]

			revealBytes, err := ioutil.ReadFile(revealFilePath)
			if err != nil {
				return err
			}

			msg := types.NewMsgRevealBid(auctionID, hex.EncodeToString(revealBytes), clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
