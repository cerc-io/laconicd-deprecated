package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cerc-io/laconicd/server/flags"
	"github.com/cerc-io/laconicd/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewTxCmd returns a root CLI command handler for all x/bond transaction commands.
func NewTxCmd() *cobra.Command {
	bondTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "nameservice transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bondTxCmd.AddCommand(
		GetCmdSetRecord(),
		GetCmdRenewRecord(),
		GetCmdAssociateBond(),
		GetCmdDissociateBond(),
		GetCmdDissociateRecords(),
		GetCmdReAssociateRecords(),
		GetCmdSetName(),
		GetCmdReserveName(),
		GetCmdSetAuthorityBond(),
		GetCmdDeleteName(),
	)

	return bondTxCmd
}

// GetCmdSetRecord is the CLI command for creating/updating a record.
func GetCmdSetRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [payload file path] [bond-id]",
		Short: "Set record.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new record with payload and bond id.
Example:
$ %s tx %s set [payload file path] [bond-id]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			payload, err := GetPayloadFromFile(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetRecord(payload.ToPayload(), args[1], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdRenewRecord is the CLI command for renewing an expired record.
func GetCmdRenewRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew-record [record-id]",
		Short: "Renew (expired) record.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Renew record.
Example:
$ %s tx %s renew-record [record-id]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgRenewRecord(args[0], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdAssociateBond is the CLI command for associating a record with a bond.
func GetCmdAssociateBond() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "associate-bond [record-id] [bond-id]",
		Short: "Associate record with bond.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Associate record with bond.
Example:
$ %s tx %s associate-bond [record-id] [bond-id]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgAssociateBond(args[0], args[1], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdDissociateBond is the CLI command for dissociating a record from a bond.
func GetCmdDissociateBond() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dissociate-bond [record-id]",
		Short: "Dissociate record from (existing) bond.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Dissociate record from (existing) bond.
Example:
$ %s tx %s dissociate-bond [record-id]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgDissociateBond(args[0], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdDissociateRecords is the CLI command for dissociating all records from a bond.
func GetCmdDissociateRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dissociate-records [bond-id]",
		Short: "Dissociate all records from bond.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Dissociate all records from bond.
Example:
$ %s tx %s dissociate-bond [record-id]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgDissociateRecords(args[0], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdReAssociateRecords is the CLI command for reassociating all records from old to new bond.
func GetCmdReAssociateRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reassociate-records [old-bond-id] [new-bond-id]",
		Short: "Re-Associates all records from old to new bond.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Re-Associates all records from old to new bond.
Example:
$ %s tx %s reassociate-records [old-bond-id] [new-bond-id]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgReAssociateRecords(args[0], args[1], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdSetName is the CLI command for mapping a name to a CID.
func GetCmdSetName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-name [crn] [cid]",
		Short: "Set CRN to CID mapping.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Set name with crn and cid.
Example:
$ %s tx %s set-name [crn] [cid]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetName(args[0], args[1], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdReserveName is the CLI command for reserving a name.
func GetCmdReserveName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reserve-name [name]",
		Short: "Reserve name.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Reserver name with owner address .
Example:
$ %s tx %s reserve-name [name] --owner [ownerAddress]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			owner, err := cmd.Flags().GetString("owner")
			if err != nil {
				return err
			}
			ownerAddress, err := sdk.AccAddressFromBech32(owner)
			if err != nil {
				return err
			}

			msg := types.NewMsgReserveAuthority(args[0], clientCtx.GetFromAddress(), ownerAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String("owner", "", "Owner address, if creating a sub-authority.")

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdSetAuthorityBond is the CLI command for associating a bond with an authority.
func GetCmdSetAuthorityBond() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authority-bond [name] [bond-id]",
		Short: "Associate authority with bond.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Reserver name with owner address .
Example:
$ %s tx %s authority-bond [name] [bond-id]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgSetAuthorityBond(args[0], args[1], clientCtx.GetFromAddress())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

func GetCmdDeleteName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-name [crn]",
		Short: "Delete CRN.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delete CRN.
Example:
$ %s tx %s delete-name [crn]
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteNameAuthority(args[0], clientCtx.GetFromAddress())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd, _ = flags.AddTxFlags(cmd)
	return cmd
}

// GetPayloadFromFile  Load payload object from YAML file.
func GetPayloadFromFile(filePath string) (*types.PayloadType, error) {
	var payload types.PayloadType

	data, err := ioutil.ReadFile(filePath) // #nosec G304
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
