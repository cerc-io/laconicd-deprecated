package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/tharsis/ethermint/x/bond/types"
	"strings"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	bondQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bond module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bondQueryCmd.AddCommand(
		GetQueryParamsCmd(),
		GetQueryBondLists(),
		GetBondByIdCmd(),
		GetBondListByOwnerCmd(),
		GetBondModuleBalanceCmd(),
	)

	return bondQueryCmd
}

// GetQueryParamsCmd implements the params query command.
func GetQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current bond parameters information.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as bond parameters.

Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetQueryBondLists implements the bond lists query command.
func GetQueryBondLists() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List bonds.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get bond list .

Example:
$ %s query %s list
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Bonds(cmd.Context(), &types.QueryGetBondsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetBondByIdCmd implements the bond info by id query command.
func GetBondByIdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [bond Id]",
		Short: "Get bond.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get bond info by bond id .

Example:
$ %s query bond get {BOND ID}
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			id := args[0]

			res, err := queryClient.GetBondById(cmd.Context(), &types.QueryGetBondByIdRequest{Id: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetBondListByOwnerCmd queries the bond list by owner.
func GetBondListByOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-by-owner [address]",
		Short: "Query bonds by owner.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get bond list by owner.

Example:
$ %s query %s query-by-owner [address]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			owner := args[0]
			res, err := queryClient.GetBondsByOwner(cmd.Context(), &types.QueryGetBondsByOwnerRequest{Owner: owner})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetBondModuleBalanceCmd queries the bond module account balance.
func GetBondModuleBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Get bond module account balance.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get bond module balance.

Example:
$ %s query %s balance
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.GetBondsModuleBalance(cmd.Context(), &types.QueryGetBondModuleBalanceRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
