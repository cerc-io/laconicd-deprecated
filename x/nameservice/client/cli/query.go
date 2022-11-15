package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cerc-io/laconicd/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	nameserviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the nameservice module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	nameserviceQueryCmd.AddCommand(
		GetCmdWhoIs(),
		GetCmdResolve(),
		GetCmdLookupCRN(),
		GetRecordExpiryQueue(),
		GetAuthorityExpiryQueue(),
		GetQueryParamsCmd(),
		GetCmdList(),
		GetCmdGetResource(),
		GetCmdQueryByBond(),
		GetCmdBalance(),
		GetCmdNames(),
	)
	return nameserviceQueryCmd
}

// GetCmdWhoIs queries a whois info for a name.
func GetCmdWhoIs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whois [name]",
		Short: "Get name owner info.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get name owner info.
Example:
$ %s query %s whois [name]
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
			res, err := queryClient.Whois(cmd.Context(), &types.QueryWhoisRequest{Name: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdLookupCRN queries naming info for a CRN.
func GetCmdLookupCRN() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookup [crn]",
		Short: "Get naming info for CRN.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get naming info for CRN.
Example:
$ %s query %s lookup [crn]
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
			res, err := queryClient.LookupCrn(cmd.Context(), &types.QueryLookupCrn{Crn: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
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
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdList queries all records.
func GetCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List records.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the records.
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
			res, err := queryClient.ListRecords(cmd.Context(), &types.QueryListRecordsRequest{})
			if err != nil {
				return err
			}

			recordsList := res.GetRecords()
			records := make([]types.RecordType, len(recordsList))
			for i, record := range res.GetRecords() {
				records[i] = record.ToRecordType()
			}
			bytesResult, err := json.Marshal(records)
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(bytesResult)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetResource queries a record record.
func GetCmdGetResource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [ID]",
		Short: "Get record.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the record by id.
Example:
$ %s query %s get [ID]
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
			record, err := queryClient.GetRecord(cmd.Context(), &types.QueryRecordByIDRequest{Id: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(record)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdResolve resolves a CRN to a record.
func GetCmdResolve() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve [crn]",
		Short: "Resolve CRN to record.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Resolve CRN to record.
Example:
$ %s query %s resolve [crn]
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
			record, err := queryClient.ResolveCrn(cmd.Context(), &types.QueryResolveCrn{Crn: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(record)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryByBond queries records by bond ID.
func GetCmdQueryByBond() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-by-bond [bond-id]",
		Short: "Query records by bond ID.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the record by bond id.
Example:
$ %s query %s query-by-bond [bond id]
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

			bondID := args[0]
			res, err := queryClient.GetRecordByBondID(cmd.Context(), &types.QueryRecordByBondIDRequest{Id: bondID})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdBalance queries the bond module account balance.
func GetCmdBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Get record rent module account balance.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the record rent module account balance.
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
			res, err := queryClient.GetNameServiceModuleBalance(cmd.Context(), &types.GetNameServiceModuleBalanceRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdNames queries all naming records.
func GetCmdNames() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "names",
		Short: "List name records.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the names list.
Example:
$ %s query %s names
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
			res, err := queryClient.ListNameRecords(cmd.Context(), &types.QueryListNameRecordsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetRecordExpiryQueue gets the record expiry queue.
func GetRecordExpiryQueue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-expiry",
		Short: "Get record expiry queue.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get record expiry queue.
Example:
$ %s query %s record-expiry
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
			res, err := queryClient.GetRecordExpiryQueue(cmd.Context(), &types.QueryGetRecordExpiryQueue{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetAuthorityExpiryQueue gets the authority expiry queue.
func GetAuthorityExpiryQueue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authority-expiry",
		Short: "Get authority expiry queue.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get authority expiry queue.
Example:
$ %s query %s authority-expiry
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
			res, err := queryClient.GetAuthorityExpiryQueue(cmd.Context(), &types.QueryGetAuthorityExpiryQueue{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
