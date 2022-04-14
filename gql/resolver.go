package gql

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	auctiontypes "github.com/tharsis/ethermint/x/auction/types"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
	nstypes "github.com/tharsis/ethermint/x/nameservice/types"
)

// DefaultLogNumLines is the number of log lines to tail by default.
const DefaultLogNumLines = 50

// MaxLogNumLines is the max number of log lines that can be tailed.
const MaxLogNumLines = 1000

type Resolver struct {
	ctx     client.Context
	logFile string
}

// Query is the entry point to query execution.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (q queryResolver) LookupAuthorities(ctx context.Context, names []string) ([]*AuthorityRecord, error) {
	nsQueryClient := nstypes.NewQueryClient(q.ctx)
	auctionQueryClient := auctiontypes.NewQueryClient(q.ctx)
	var gqlResponse []*AuthorityRecord

	for _, name := range names {
		res, err := nsQueryClient.Whois(context.Background(), &nstypes.QueryWhoisRequest{Name: name})
		if err != nil {
			return nil, err
		}

		nameAuthority := res.GetNameAuthority()
		gqlNameAuthorityRecord, err := GetGQLNameAuthorityRecord(&nameAuthority)
		if err != nil {
			return nil, err
		}

		if nameAuthority.AuctionId != "" {
			auctionResp, err := auctionQueryClient.GetAuction(context.Background(), &auctiontypes.AuctionRequest{Id: nameAuthority.GetAuctionId()})
			if err != nil {
				return nil, err
			}
			bidsResp, err := auctionQueryClient.GetBids(context.Background(), &auctiontypes.BidsRequest{AuctionId: nameAuthority.GetAuctionId()})
			if err != nil {
				return nil, err
			}

			gqlAuctionRecord, err := GetGQLAuction(auctionResp.GetAuction(), bidsResp.GetBids())
			if err != nil {
				return nil, err
			}

			gqlNameAuthorityRecord.Auction = gqlAuctionRecord
		}

		gqlResponse = append(gqlResponse, gqlNameAuthorityRecord)
	}

	return gqlResponse, nil
}

func (q queryResolver) ResolveNames(ctx context.Context, names []string) ([]*Record, error) {
	nsQueryClient := nstypes.NewQueryClient(q.ctx)
	var gqlResponse []*Record
	for _, name := range names {
		res, err := nsQueryClient.ResolveWrn(context.Background(), &nstypes.QueryResolveWrn{Wrn: name})
		if err != nil {
			// Return nil for record not found.
			gqlResponse = append(gqlResponse, nil)
		} else {
			gqlRecord, err := getGQLRecord(context.Background(), q, *res.GetRecord())
			if err != nil {
				return nil, err
			}

			gqlResponse = append(gqlResponse, gqlRecord)
		}
	}

	return gqlResponse, nil
}

func (q queryResolver) LookupNames(ctx context.Context, names []string) ([]*NameRecord, error) {
	nsQueryClient := nstypes.NewQueryClient(q.ctx)
	var gqlResponse []*NameRecord

	for _, name := range names {
		res, err := nsQueryClient.LookupWrn(context.Background(), &nstypes.QueryLookupWrn{Wrn: name})
		if err != nil {
			// Return nil for name not found.
			gqlResponse = append(gqlResponse, nil)
		} else {
			gqlRecord, err := getGQLNameRecord(res.GetName())
			if err != nil {
				return nil, err
			}

			gqlResponse = append(gqlResponse, gqlRecord)
		}
	}

	return gqlResponse, nil
}

func (q queryResolver) QueryRecords(ctx context.Context, attributes []*KeyValueInput, all *bool) ([]*Record, error) {
	nsQueryClient := nstypes.NewQueryClient(q.ctx)

	res, err := nsQueryClient.ListRecords(
		context.Background(),
		&nstypes.QueryListRecordsRequest{
			Attributes: parseRequestAttributes(attributes),
			All:        (all != nil && *all),
		},
	)

	if err != nil {
		return nil, err
	}

	records := res.GetRecords()
	gqlResponse := make([]*Record, len(records))

	for i, record := range records {
		gqlRecord, err := getGQLRecord(context.Background(), q, record)
		if err != nil {
			return nil, err
		}
		gqlResponse[i] = gqlRecord
	}

	return gqlResponse, nil

}

func (q queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*Record, error) {
	nsQueryClient := nstypes.NewQueryClient(q.ctx)
	gqlResponse := make([]*Record, len(ids))

	for i, id := range ids {
		res, err := nsQueryClient.GetRecord(context.Background(), &nstypes.QueryRecordByIdRequest{Id: id})
		if err != nil {
			// Return nil for record not found.
			gqlResponse[i] = nil
		} else {
			record, err := getGQLRecord(context.Background(), q, res.GetRecord())
			if err != nil {
				return nil, err
			}
			gqlResponse[i] = record
		}
	}

	return gqlResponse, nil
}

func (q queryResolver) GetStatus(ctx context.Context) (*Status, error) {
	nodeInfo, syncInfo, validatorInfo, err := getStatusInfo(q.ctx)
	if err != nil {
		return nil, err
	}

	numPeers, peers, err := getNetInfo(q.ctx)
	if err != nil {
		return nil, err
	}

	validatorSet, err := getValidatorSet(q.ctx)
	if err != nil {
		return nil, err
	}

	diskUsage, err := GetDiskUsage(NodeDataPath)
	if err != nil {
		return nil, err
	}

	return &Status{
		Version:    NameServiceVersion,
		Node:       nodeInfo,
		Sync:       syncInfo,
		Validator:  validatorInfo,
		Validators: validatorSet,
		NumPeers:   numPeers,
		Peers:      peers,
		DiskUsage:  diskUsage,
	}, nil
}

func (q queryResolver) GetAccounts(ctx context.Context, addresses []string) ([]*Account, error) {
	accounts := make([]*Account, len(addresses))
	for index, address := range addresses {
		account, err := q.GetAccount(ctx, address)
		if err != nil {
			return nil, err
		}
		accounts[index] = account
	}
	return accounts, nil
}

func (q queryResolver) GetAccount(ctx context.Context, address string) (*Account, error) {
	authQueryClient := authtypes.NewQueryClient(q.ctx)
	accountResponse, err := authQueryClient.Account(ctx, &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		return nil, err
	}
	var account authtypes.AccountI
	err = q.ctx.Codec.UnpackAny(accountResponse.GetAccount(), &account)
	if err != nil {
		return nil, err
	}
	var pubKey *string
	if account.GetPubKey() != nil {
		pubKeyStr := base64.StdEncoding.EncodeToString(account.GetPubKey().Bytes())
		pubKey = &pubKeyStr
	}

	// Get the account balance
	bankQueryClient := banktypes.NewQueryClient(q.ctx)
	balance, err := bankQueryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: address})

	accNum := strconv.FormatUint(account.GetAccountNumber(), 10)
	seq := strconv.FormatUint(account.GetSequence(), 10)

	return &Account{
		Address:  address,
		Number:   accNum,
		Sequence: seq,
		PubKey:   pubKey,
		Balance:  getGQLCoins(balance.GetBalances()),
	}, nil
}

func (q queryResolver) GetBondsByIds(ctx context.Context, ids []string) ([]*Bond, error) {
	bonds := make([]*Bond, len(ids))
	for index, id := range ids {
		bondObj, err := q.GetBond(ctx, id)
		if err != nil {
			return nil, err
		}
		bonds[index] = bondObj
	}

	return bonds, nil
}

func (q *queryResolver) GetBond(ctx context.Context, id string) (*Bond, error) {
	bondQueryClient := bondtypes.NewQueryClient(q.ctx)
	bondResp, err := bondQueryClient.GetBondById(context.Background(), &bondtypes.QueryGetBondByIdRequest{Id: id})
	if err != nil {
		return nil, err
	}

	bond := bondResp.GetBond()
	if bond == nil {
		return nil, nil
	}
	return getGQLBond(bondResp.GetBond())
}

func (q queryResolver) QueryBonds(ctx context.Context, attributes []*KeyValueInput) ([]*Bond, error) {
	bondQueryClient := bondtypes.NewQueryClient(q.ctx)
	bonds, err := bondQueryClient.Bonds(context.Background(), &bondtypes.QueryGetBondsRequest{})
	if err != nil {
		return nil, err
	}

	gqlResponse := make([]*Bond, len(bonds.GetBonds()))
	for i, bondObj := range bonds.GetBonds() {
		gqlBond, err := getGQLBond(bondObj)
		if err != nil {
			return nil, err
		}
		gqlResponse[i] = gqlBond
	}

	return gqlResponse, nil
}

// QueryBondsByOwner will return bonds by owner
func (q queryResolver) QueryBondsByOwner(ctx context.Context, ownerAddresses []string) ([]*OwnerBonds, error) {
	ownerBonds := make([]*OwnerBonds, len(ownerAddresses))
	for index, ownerAddress := range ownerAddresses {
		bondsObj, err := q.GetBondsByOwner(ctx, ownerAddress)
		if err != nil {
			return nil, err
		}
		ownerBonds[index] = bondsObj
	}

	return ownerBonds, nil
}

func (q queryResolver) GetBondsByOwner(ctx context.Context, address string) (*OwnerBonds, error) {
	bondQueryClient := bondtypes.NewQueryClient(q.ctx)
	bondResp, err := bondQueryClient.GetBondsByOwner(context.Background(), &bondtypes.QueryGetBondsByOwnerRequest{Owner: address})
	if err != nil {
		return nil, err
	}

	ownerBonds := make([]*Bond, len(bondResp.GetBonds()))
	for i, bond := range bondResp.GetBonds() {
		bondObj, err := getGQLBond(&bond)
		if err != nil {
			return nil, err
		}
		ownerBonds[i] = bondObj
	}

	return &OwnerBonds{Bonds: ownerBonds, Owner: address}, nil
}

func (q queryResolver) GetAuctionsByIds(ctx context.Context, ids []string) ([]*Auction, error) {
	auctionQueryClient := auctiontypes.NewQueryClient(q.ctx)
	gqlAuctionResponse := make([]*Auction, len(ids))
	for i, id := range ids {
		auctionObj, err := auctionQueryClient.GetAuction(context.Background(), &auctiontypes.AuctionRequest{Id: id})
		if err != nil {
			return nil, err
		}
		bidsObj, err := auctionQueryClient.GetBids(context.Background(), &auctiontypes.BidsRequest{AuctionId: id})
		if err != nil {
			return nil, err
		}

		gqlAuction, err := GetGQLAuction(auctionObj.GetAuction(), bidsObj.GetBids())
		if err != nil {
			return nil, err
		}

		gqlAuctionResponse[i] = gqlAuction
	}

	return gqlAuctionResponse, nil
}
