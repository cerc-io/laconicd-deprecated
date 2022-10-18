package nameservice

import (
	"time"

	"github.com/cerc-io/laconicd/x/nameservice/keeper"
	"github.com/cerc-io/laconicd/x/nameservice/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)

	for _, record := range data.Records {
		keeper.PutRecord(ctx, record)

		// Add to record expiry queue if expiry time is in the future.
		expiryTime, err := time.Parse(time.RFC3339, record.ExpiryTime)
		if err != nil {
			panic(err)
		}

		if expiryTime.After(ctx.BlockTime()) {
			keeper.InsertRecordExpiryQueue(ctx, record)
		}

		// Note: Bond genesis runs first, so bonds will already be present.
		if record.BondId != "" {
			keeper.AddBondToRecordIndexEntry(ctx, record.BondId, record.Id)
		}
	}

	for _, authority := range data.Authorities {
		// Only import authorities that are marked active.
		if authority.Entry.Status == types.AuthorityActive {
			keeper.SetNameAuthority(ctx, authority.Name, authority.Entry)

			// Add authority name to expiry queue.
			keeper.InsertAuthorityExpiryQueue(ctx, authority.Name, authority.Entry.ExpiryTime)

			// Note: Bond genesis runs first, so bonds will already be present.
			if authority.Entry.BondId != "" {
				keeper.AddBondToAuthorityIndexEntry(ctx, authority.Entry.BondId, authority.Name)
			}
		}
	}

	for _, nameEntry := range data.Names {
		keeper.SetNameRecord(ctx, nameEntry.Name, nameEntry.Entry.Latest.Id)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)

	records := keeper.ListRecords(ctx)

	authorities := keeper.ListNameAuthorityRecords(ctx)
	authorityEntries := []types.AuthorityEntry{}
	// #nosec
	for name, record := range authorities {
		authorityEntries = append(authorityEntries, types.AuthorityEntry{
			Name:  name,
			Entry: &record, //nolint: all
		})
	}

	names := keeper.ListNameRecords(ctx)

	return types.GenesisState{
		Params:      params,
		Records:     records,
		Authorities: authorityEntries,
		Names:       names,
	}
}
