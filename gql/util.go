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

func parseRequestValue(value *ValueInput) *registrytypes.QueryListRecordsRequest_ValueInput {
	if value == nil {
		return nil
	}
	var val registrytypes.QueryListRecordsRequest_ValueInput

	if value.String != nil {
		val.String_ = *value.String
		val.Type = "string"
	}

	if value.Int != nil {
		val.Int = int64(*value.Int)
		val.Type = "int"
	}

	if value.Float != nil {
		val.Float = *value.Float
		val.Type = "float"
	}

	if value.Boolean != nil {
		val.Boolean = *value.Boolean
		val.Type = "boolean"
	}

	if value.Link != nil {
		reference := &registrytypes.QueryListRecordsRequest_ReferenceInput{
			Id: value.Link.String(),
		}

		val.Reference = reference
		val.Type = "reference"
	}

	// handle arrays
	if value.Array != nil {
		values := []*registrytypes.QueryListRecordsRequest_ValueInput{}
		for _, v := range value.Array {
			val := parseRequestValue(v)
			values = append(values, val)
		}
		val.Values = values
		val.Type = "array"
	}

	return &val
}

func parseRequestAttributes(attrs []*KeyValueInput) []*registrytypes.QueryListRecordsRequest_KeyValueInput {
	kvPairs := []*registrytypes.QueryListRecordsRequest_KeyValueInput{}

	for _, value := range attrs {
		parsedValue := parseRequestValue(value.Value)
		kvPair := &registrytypes.QueryListRecordsRequest_KeyValueInput{
			Key:   value.Key,
			Value: parsedValue,
		}
		kvPairs = append(kvPairs, kvPair)
	}

	return kvPairs
}
