package gql

import (
	"context"
	"fmt" // #nosec G702
	"strconv"

	auctiontypes "github.com/cerc-io/laconicd/x/auction/types"
	bondtypes "github.com/cerc-io/laconicd/x/bond/types"
	registrytypes "github.com/cerc-io/laconicd/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
)

// OwnerAttributeName denotes the owner attribute name for a bond.
const OwnerAttributeName = "owner"

// BondIDAttributeName denotes the record bond ID.
const BondIDAttributeName = "bondId"

// ExpiryTimeAttributeName denotes the record expiry time.
const ExpiryTimeAttributeName = "expiryTime"

func getGQLCoin(coin sdk.Coin) *Coin {
	gqlCoin := Coin{
		Type:     coin.Denom,
		Quantity: coin.Amount.BigInt().String(),
	}

	return &gqlCoin
}

func getGQLCoins(coins sdk.Coins) []*Coin {
	gqlCoins := make([]*Coin, len(coins))
	for index, coin := range coins {
		gqlCoins[index] = getGQLCoin(coin)
	}

	return gqlCoins
}

func GetGQLNameAuthorityRecord(record *registrytypes.NameAuthority) (*AuthorityRecord, error) {
	if record == nil {
		return nil, nil
	}

	return &AuthorityRecord{
		OwnerAddress:   record.OwnerAddress,
		OwnerPublicKey: record.OwnerPublicKey,
		Height:         strconv.FormatUint(record.Height, 10),
		Status:         record.Status,
		BondID:         record.GetBondId(),
		ExpiryTime:     record.GetExpiryTime().String(),
	}, nil
}

func getGQLRecord(ctx context.Context, resolver QueryResolver, record registrytypes.Record) (*Record, error) {
	// Nil record.
	if record.Deleted {
		return nil, nil
	}

	node, err := ipld.Decode(record.Attributes, dagjson.Decode)
	if err != nil {
		return nil, err
	}
	if node.Kind() != ipld.Kind_Map {
		return nil, fmt.Errorf("invalid record attributes")
	}

	var links []string
	attributes, err := resolveIPLDNode(node, &links)
	if err != nil {
		return nil, err
	}

	references, err := resolver.GetRecordsByIds(ctx, links)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:         record.Id,
		BondID:     record.GetBondId(),
		CreateTime: record.GetCreateTime(),
		ExpiryTime: record.GetExpiryTime(),
		Owners:     record.GetOwners(),
		Names:      record.GetNames(),
		Attributes: attributes.(MapValue).Value,
		References: references,
	}, nil
}

func resolveIPLDNode(node ipld.Node, links *[]string) (Value, error) {
	switch node.Kind() {
	case ipld.Kind_Map:
		var entries []*Attribute
		for itr := node.MapIterator(); !itr.Done(); {
			k, v, err := itr.Next()
			if err != nil {
				return nil, err
			}
			if k.Kind() != ipld.Kind_String {
				return nil, fmt.Errorf("invalid record attribute key type: %s", k.Kind())
			}
			s, err := k.AsString()
			if err != nil {
				return nil, err
			}
			val, err := resolveIPLDNode(v, links)
			if err != nil {
				return nil, err
			}
			entries = append(entries, &Attribute{
				Key:   s,
				Value: val,
			})
		}
		return MapValue{entries}, nil
	case ipld.Kind_List:
		var values []Value
		for itr := node.ListIterator(); !itr.Done(); {
			_, v, err := itr.Next()
			if err != nil {
				return nil, err
			}
			val, err := resolveIPLDNode(v, links)
			if err != nil {
				return nil, err
			}
			values = append(values, val)
		}
		return ArrayValue{values}, nil
	case ipld.Kind_Null:
		return nil, nil
	case ipld.Kind_Bool:
		val, err := node.AsBool()
		if err != nil {
			return nil, err
		}
		return BooleanValue{val}, nil
	case ipld.Kind_Int:
		val, err := node.AsInt()
		if err != nil {
			return nil, err
		}
		// TODO: handle bigger ints
		return IntValue{int(val)}, nil
	case ipld.Kind_Float:
		val, err := node.AsFloat()
		if err != nil {
			return nil, err
		}
		return FloatValue{val}, nil
	case ipld.Kind_String:
		val, err := node.AsString()
		if err != nil {
			return nil, err
		}
		return StringValue{val}, nil
	case ipld.Kind_Bytes:
		val, err := node.AsBytes()
		if err != nil {
			return nil, err
		}
		return BytesValue{string(val)}, nil
	case ipld.Kind_Link:
		val, err := node.AsLink()
		if err != nil {
			return nil, err
		}
		*links = append(*links, val.String())
		return LinkValue{Link(val.String())}, nil
	default:
		return nil, fmt.Errorf("invalid node kind")
	}
}

func getGQLNameRecord(record *registrytypes.NameRecord) (*NameRecord, error) {
	if record == nil {
		return nil, fmt.Errorf("got nil record")
	}

	records := make([]*NameRecordEntry, len(record.History))
	for index, entry := range record.History {
		records[index] = getNameRecordEntry(entry)
	}

	return &NameRecord{
		Latest:  getNameRecordEntry(record.Latest),
		History: records,
	}, nil
}

func getNameRecordEntry(record *registrytypes.NameRecordEntry) *NameRecordEntry {
	return &NameRecordEntry{
		ID:     record.Id,
		Height: strconv.FormatUint(record.Height, 10),
	}
}

func getGQLBond(bondObj *bondtypes.Bond) (*Bond, error) {
	// Nil record.
	if bondObj == nil {
		return nil, nil
	}

	return &Bond{
		ID:      bondObj.Id,
		Owner:   bondObj.Owner,
		Balance: getGQLCoins(bondObj.Balance),
	}, nil
}

func getAuctionBid(bid *auctiontypes.Bid) *AuctionBid {
	return &AuctionBid{
		BidderAddress: bid.BidderAddress,
		Status:        bid.Status,
		CommitHash:    bid.CommitHash,
		CommitTime:    bid.GetCommitTime(),
		RevealTime:    bid.GetRevealTime(),
		CommitFee:     getGQLCoin(bid.CommitFee),
		RevealFee:     getGQLCoin(bid.RevealFee),
		BidAmount:     getGQLCoin(bid.BidAmount),
	}
}

func GetGQLAuction(auction *auctiontypes.Auction, bids []*auctiontypes.Bid) (*Auction, error) {
	if auction == nil {
		return nil, nil
	}

	gqlAuction := Auction{
		ID:             auction.Id,
		Status:         auction.Status,
		OwnerAddress:   auction.OwnerAddress,
		CreateTime:     auction.GetCreateTime(),
		CommitsEndTime: auction.GetCommitsEndTime(),
		RevealsEndTime: auction.GetRevealsEndTime(),
		CommitFee:      getGQLCoin(auction.CommitFee),
		RevealFee:      getGQLCoin(auction.RevealFee),
		MinimumBid:     getGQLCoin(auction.MinimumBid),
		WinnerAddress:  auction.WinnerAddress,
		WinnerBid:      getGQLCoin(auction.WinningBid),
		WinnerPrice:    getGQLCoin(auction.WinningPrice),
	}

	auctionBids := make([]*AuctionBid, len(bids))
	for index, entry := range bids {
		auctionBids[index] = getAuctionBid(entry)
	}

	gqlAuction.Bids = auctionBids

	return &gqlAuction, nil
}

func toRPCValue(value *ValueInput) *registrytypes.QueryListRecordsRequest_ValueInput {
	var rpcval registrytypes.QueryListRecordsRequest_ValueInput

	switch {
	case value == nil:
		return nil
	case value.Int != nil:
		rpcval.Value = &registrytypes.QueryListRecordsRequest_ValueInput_Int{Int: int64(*value.Int)}
	case value.Float != nil:
		rpcval.Value = &registrytypes.QueryListRecordsRequest_ValueInput_Float{Float: *value.Float}
	case value.String != nil:
		rpcval.Value = &registrytypes.QueryListRecordsRequest_ValueInput_String_{String_: *value.String}
	case value.Boolean != nil:
		rpcval.Value = &registrytypes.QueryListRecordsRequest_ValueInput_Boolean{Boolean: *value.Boolean}
	case value.Link != nil:
		rpcval.Value = &registrytypes.QueryListRecordsRequest_ValueInput_Link{Link: value.Link.String()}
	case value.Array != nil:
		var contents registrytypes.QueryListRecordsRequest_ArrayInput
		for _, val := range value.Array {
			contents.Values = append(contents.Values, toRPCValue(val))
		}
		rpcval.Value = &registrytypes.QueryListRecordsRequest_ValueInput_Array{Array: &contents}
	case value.Map != nil:
		var contents registrytypes.QueryListRecordsRequest_MapInput
		for _, kv := range value.Map {
			contents.Values[kv.Key] = toRPCValue(kv.Value)
		}
		rpcval.Value = &registrytypes.QueryListRecordsRequest_ValueInput_Map{Map: &contents}
	}
	return &rpcval
}

func toRPCAttributes(attrs []*KeyValueInput) []*registrytypes.QueryListRecordsRequest_KeyValueInput {
	kvPairs := []*registrytypes.QueryListRecordsRequest_KeyValueInput{}

	for _, value := range attrs {
		parsedValue := toRPCValue(value.Value)
		kvPair := &registrytypes.QueryListRecordsRequest_KeyValueInput{
			Key:   value.Key,
			Value: parsedValue,
		}
		kvPairs = append(kvPairs, kvPair)
	}

	return kvPairs
}
