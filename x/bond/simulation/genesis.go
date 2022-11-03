package simulation

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cerc-io/laconicd/x/bond/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RandomizedGenState generates a random GenesisState
func RandomizedGenState(simState *module.SimulationState) {
	bondParams := types.NewParams(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))))
	bondGenesis := types.NewGenesisState(bondParams, []*types.Bond{})

	bz, err := json.MarshalIndent(bondGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(bondGenesis)
}
