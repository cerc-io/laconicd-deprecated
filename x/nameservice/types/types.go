package types

import (
	"crypto/sha256"
	canonicalJson "github.com/gibson042/canonicaljson-go"
	"github.com/tharsis/ethermint/x/nameservice/helpers"
	"time"
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
func (payloadObj *PayloadType) ToPayload() Payload {
	var payload = Payload{
		Record: &Record{
			CreateTime: time.Time{},
			ExpiryTime: time.Time{},
			Deleted:    false,
			Owners:     nil,
			Attributes: helpers.MarshalMapToJSONBytes(payloadObj.Record),
		},
		Signatures: payloadObj.Signatures,
	}
	return payload
}

// ToReadablePayload converts Payload to PayloadType
// It will unmarshal with record attributes
func (payload Payload) ToReadablePayload() PayloadType {
	var payloadType PayloadType

	payloadType.Record = helpers.UnMarshalMapFromJSONBytes(payload.Record.Attributes)

	payloadType.Signatures = payload.Signatures

	return payloadType
}

// Record to Record Type for human-readable attributes

func (r *Record) ToRecordType() RecordType {
	var resourceObj RecordType

	resourceObj.Id = r.Id
	resourceObj.BondId = r.BondId
	resourceObj.CreateTime = r.CreateTime
	resourceObj.ExpiryTime = r.ExpiryTime
	resourceObj.Deleted = r.Deleted
	resourceObj.Owners = r.Owners
	resourceObj.Attributes = helpers.UnMarshalMapFromJSONBytes(r.Attributes)

	return resourceObj
}

// RecordType represents a WNS record.
type RecordType struct {
	Id         string                 `json:"id,omitempty"`
	Names      []string               `json:"names,omitempty"`
	BondId     string                 `json:"bondId,omitempty"`
	CreateTime time.Time              `json:"createTime,omitempty"`
	ExpiryTime time.Time              `json:"expiryTime,omitempty"`
	Deleted    bool                   `json:"deleted,omitempty"`
	Owners     []string               `json:"owners,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// ToRecordObj converts Record to RecordObj.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (r *RecordType) ToRecordObj() Record {
	var resourceObj Record

	resourceObj.Id = r.Id
	resourceObj.BondId = r.BondId
	resourceObj.CreateTime = r.CreateTime
	resourceObj.ExpiryTime = r.ExpiryTime
	resourceObj.Deleted = r.Deleted
	resourceObj.Owners = r.Owners
	resourceObj.Attributes = helpers.MarshalMapToJSONBytes(r.Attributes)

	return resourceObj
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
