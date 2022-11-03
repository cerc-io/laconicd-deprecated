package simulation

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cerc-io/laconicd/x/nameservice/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RandomizedGenState generates a random GenesisState
func RandomizedGenState(simState *module.SimulationState) {
	nameserviceParams := types.NewParams(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
		time.Duration(simState.Rand.Intn(1000))*time.Second,
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
		time.Duration(simState.Rand.Intn(1000))*time.Second,
		time.Duration(simState.Rand.Intn(1000))*time.Second,
		false,
		time.Duration(simState.Rand.Intn(1000))*time.Second,
		time.Duration(simState.Rand.Intn(1000))*time.Second,
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
	)

	nameserviceGenesis := types.NewGenesisState(nameserviceParams, []types.Record{}, []types.AuthorityEntry{}, []types.NameEntry{})

	bz, err := json.MarshalIndent(nameserviceGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&nameserviceGenesis)
}
