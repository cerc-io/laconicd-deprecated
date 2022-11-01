package simulation

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cerc-io/laconicd/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RandomizedGenState generates a random GenesisState
func RandomizedGenState(simState *module.SimulationState) {
	auctionParams := types.NewParams(time.Duration(simState.Rand.Intn(1000))*time.Second,
		time.Duration(simState.Rand.Intn(1000))*time.Second,
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simState.Rand.Intn(10000000000)))),
	)

	auctionGenesis := types.NewGenesisState(auctionParams, []*types.Auction{})

	bz, err := json.MarshalIndent(auctionGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(auctionGenesis)
}
