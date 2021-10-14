package flags

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

// Tendermint full-node start flags
const (
	WithTendermint = "with-tendermint"
	Address        = "address"
	Transport      = "transport"
	TraceStore     = "trace-store"
	CPUProfile     = "cpu-profile"
)

// GRPC-related flags.
const (
	GRPCEnable     = "grpc.enable"
	GRPCAddress    = "grpc.address"
	GRPCWebEnable  = "grpc-web.enable"
	GRPCWebAddress = "grpc-web.address"
)

// JSON-RPC flags
const (
	JSONRPCEnable  = "json-rpc.enable"
	JSONRPCAPI     = "json-rpc.api"
	JSONRPCAddress = "json-rpc.address"
	JSONWsAddress  = "json-rpc.ws-address"
	JSONRPCGasCap  = "json-rpc.gas-cap"
)

// EVM flags
const (
	EVMTracer = "evm.tracer"
)

// AddTxFlags adds common flags for commands to post tx
func AddTxFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().String(flags.FlagChainID, "testnet", "Specify Chain ID for sending Tx")
	cmd.PersistentFlags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	cmd.PersistentFlags().String(flags.FlagFees, "", "Fees to pay along with transaction; eg: 10aphoton")
	cmd.PersistentFlags().String(flags.FlagGasPrices, "", "Gas prices to determine the transaction fee (e.g. 10aphoton)")
	cmd.PersistentFlags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	cmd.PersistentFlags().Float64(flags.FlagGasAdjustment, flags.DefaultGasAdjustment, "adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ")
	cmd.PersistentFlags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, keyring.BackendOS, "Select keyring's backend")
	cmd.PersistentFlags().BoolP(flags.FlagSkipConfirmation, "y", false, "Skip tx broadcasting prompt confirmation")
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")

	// --gas can accept integers and "simulate"
	// cmd.PersistentFlags().Var(&flags.GasFlagVar, "gas", fmt.Sprintf(
	//	"gas limit to set per-transaction; set to %q to calculate required gas automatically (default %d)",
	//	flags.GasFlagAuto, flags.DefaultGasLimit,
	// ))

	// viper.BindPFlag(flags.FlagTrustNode, cmd.Flags().Lookup(flags.FlagTrustNode))

	// TODO: we need to handle the errors for these, decide if we should return error upward and handle
	// nolint: errcheck
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	// nolint: errcheck
	viper.BindPFlag(flags.FlagKeyringBackend, cmd.Flags().Lookup(flags.FlagKeyringBackend))
	// nolint: errcheck
	cmd.MarkFlagRequired(flags.FlagChainID)
	return cmd
}

// AddGQLFlags adds gql flags for
func AddGQLFlags(cmd *cobra.Command) *cobra.Command {
	// Add flags for GQL server.
	cmd.PersistentFlags().Bool("gql-server", false, "Start GQL server.")
	cmd.PersistentFlags().Bool("gql-playground", false, "Enable GQL playground.")
	cmd.PersistentFlags().String("gql-playground-api-base", "", "GQL API base path to use in GQL playground.")
	cmd.PersistentFlags().String("gql-port", "9473", "Port to use for the GQL server.")
	cmd.PersistentFlags().String("log-file", "", "File to tail for GQL 'getLogs' API.")

	return cmd
}
