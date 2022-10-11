package cli

import (
	"github.com/cerc-io/laconicd/server/flags"

	"github.com/cerc-io/laconicd/x/bond/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/bond transaction commands.
func NewTxCmd() *cobra.Command {
	bondTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "bond transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bondTxCmd.AddCommand(
		NewCreateBondCmd(),
		RefillBondCmd(),
		WithdrawBondCmd(),
		CancelBondCmd(),
	)

	return bondTxCmd
}

// NewCreateBondCmd is the CLI command for creating a bond.
func NewCreateBondCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [amount]",
		Short: "Create bond.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateBond(sdk.NewCoins(coin), clientCtx.GetFromAddress())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlags(cmd)
	// flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// RefillBondCmd is the CLI command for creating a bond.
func RefillBondCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refill [bond Id] [amount]",
		Short: "Refill bond.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			bondId := args[0]
			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRefillBond(bondId, coin, clientCtx.GetFromAddress())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlags(cmd)
	// flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// WithdrawBondCmd is the CLI command for withdrawing funds from a bond.
func WithdrawBondCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [bond Id] [amount]",
		Short: "Withdraw amount from bond.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			bondId := args[0]
			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawBond(bondId, coin, clientCtx.GetFromAddress())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlags(cmd)
	// flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CancelBondCmd is the CLI command for cancelling a bond.
func CancelBondCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [bond Id]",
		Short: "cancel bond.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			bondId := args[0]
			msg := types.NewMsgCancelBond(bondId, clientCtx.GetFromAddress())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlags(cmd)
	// flags.AddTxFlagsToCmd(cmd)

	return cmd
}
