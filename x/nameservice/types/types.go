package types

import (
	"crypto/sha256"
	"fmt"

	"github.com/cerc-io/laconicd/x/nameservice/helpers"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	canonicalJson "github.com/gibson042/canonicaljson-go"
	"github.com/golang/protobuf/proto"
)

const (
	AuthorityActive       = "active"
	AuthorityExpired      = "expired"
	AuthorityUnderAuction = "auction"
)

// PayloadType represents a signed record payload that can be serialized from/to YAML.
type PayloadType struct {
	Record     map[string]interface{} `json:"record"`
	Signatures []Signature            `json:"signatures"`
}

// ToPayload converts PayloadType to Payload object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (payloadObj *PayloadType) ToPayload() (Payload, error) {

	attributes, err := payLoadAttributes(payloadObj.Record)
	if err != nil {
		return Payload{}, err
	}

	var payload = Payload{
		Record: &Record{
			Deleted:    false,
			Owners:     nil,
			Attributes: attributes,
		},
		Signatures: payloadObj.Signatures,
	}
	return payload, nil
}

func payLoadAttributes(recordPayLoad map[string]interface{}) (*codectypes.Any, error) {

	recordType, ok := recordPayLoad["Type"]
	if !ok {
		return &codectypes.Any{}, fmt.Errorf("cannot get type from payload")
	}
	bz := helpers.MarshalMapToJSONBytes(recordPayLoad)

	switch recordType.(string) {
	case "ServiceProviderRegistration":
		{
			var attributes ServiceProviderRegistration
			err := proto.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "WebsiteRegistrationRecord":
		{
			var attributes WebsiteRegistrationRecord
			err := proto.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	default:
		return &codectypes.Any{}, fmt.Errorf("unsupported record type %s", recordType.(string))
	}

}

// ToReadablePayload converts Payload to PayloadType
// It will unmarshal with record attributes
func (payload Payload) ToReadablePayload() PayloadType {
	var payloadType PayloadType
	payloadType.Record = helpers.UnMarshalMapFromJSONBytes(payload.Record.Attributes.Value)

	payloadType.Signatures = payload.Signatures

	return payloadType
}

// Record to Record Type for human-readable attributes

func (r *Record) ToRecordType() RecordType {
	var resourceObj RecordType

	resourceObj.ID = r.Id
	resourceObj.BondID = r.BondId
	resourceObj.CreateTime = r.CreateTime
	resourceObj.ExpiryTime = r.ExpiryTime
	resourceObj.Deleted = r.Deleted
	resourceObj.Owners = r.Owners
	resourceObj.Names = r.Names
	resourceObj.Attributes = helpers.UnMarshalMapFromJSONBytes(r.Attributes.Value)

	return resourceObj
}

// RecordType represents a WNS record.
type RecordType struct {
	ID         string                 `json:"id,omitempty"`
	Names      []string               `json:"names,omitempty"`
	BondID     string                 `json:"bondId,omitempty"`
	CreateTime string                 `json:"createTime,omitempty"`
	ExpiryTime string                 `json:"expiryTime,omitempty"`
	Deleted    bool                   `json:"deleted,omitempty"`
	Owners     []string               `json:"owners,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// ToRecordObj converts Record to RecordObj.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (r *RecordType) ToRecordObj() (Record, error) {
	attributes, err := payLoadAttributes(r.Attributes)
	if err != nil {
		return Record{}, err
	}

	var resourceObj Record

	resourceObj.Id = r.ID
	resourceObj.BondId = r.BondID
	resourceObj.CreateTime = r.CreateTime
	resourceObj.ExpiryTime = r.ExpiryTime
	resourceObj.Deleted = r.Deleted
	resourceObj.Owners = r.Owners
	resourceObj.Attributes = attributes

	return resourceObj, nil
}

// CanonicalJSON returns the canonical JSON representation of the record.
func (r *RecordType) CanonicalJSON() []byte {
	bytes, err := canonicalJson.Marshal(r.Attributes)
	if err != nil {
		panic("Record marshal error.")
	}

	return bytes
}

// GetSignBytes generates a record hash to be signed.
func (r *RecordType) GetSignBytes() ([]byte, []byte) {
	// Double SHA256 hash.

	// Input to the first round of hashing.
	bytes := r.CanonicalJSON()

	// First round.
	first := sha256.New()
	first.Write(bytes)
	firstHash := first.Sum(nil)

	// Second round of hashing takes as input the output of the first round.
	second := sha256.New()
	second.Write(firstHash)
	secondHash := second.Sum(nil)

	return secondHash, bytes
}

// GetCID gets the record CID.
func (r *RecordType) GetCID() (string, error) {
	id, err := helpers.GetCid(r.CanonicalJSON())
	if err != nil {
		return "", err
	}

	return id, nil
}
