package types

import (
	"crypto/sha256"

	"github.com/cerc-io/laconicd/x/registry/helpers"
	"github.com/gibson042/canonicaljson-go"
)

const (
	AuthorityActive       = "active"
	AuthorityExpired      = "expired"
	AuthorityUnderAuction = "auction"
)

// TODO if schema records are to be more permissive than allowing a map of fields, this type will
// become specific to content records. schema records will either occupy a new message or have new
// more general purpose helper types.

type AttributeMap map[string]interface{}

// ReadablePayload represents a signed record payload that can be serialized from/to YAML.
type ReadablePayload struct {
	RecordAttributes AttributeMap `json:"record" yaml:"record"`
	Signatures       []Signature  `json:"signatures" yaml:"signatures"`
}

// ReadableRecord represents a WNS record.
type ReadableRecord struct {
	ID         string       `json:"id,omitempty"`
	Names      []string     `json:"names,omitempty"`
	BondID     string       `json:"bondId,omitempty"`
	CreateTime string       `json:"createTime,omitempty"`
	ExpiryTime string       `json:"expiryTime,omitempty"`
	Deleted    bool         `json:"deleted,omitempty"`
	Owners     []string     `json:"owners,omitempty"`
	Attributes AttributeMap `json:"attributes,omitempty"`
}

// ToPayload converts PayloadEncodable to Payload object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (payloadObj *ReadablePayload) ToPayload() Payload {
	// Note: record directly contains the attributes here
	attributes := payloadObj.RecordAttributes
	payload := Payload{
		Record: &Record{
			Deleted:    false,
			Owners:     nil,
			Attributes: helpers.MustMarshalJSON(attributes),
		},
		Signatures: payloadObj.Signatures,
	}
	return payload
}

// ToReadablePayload converts Payload to a serializable object
func (payload Payload) ToReadablePayload() ReadablePayload {
	var encodable ReadablePayload

	encodable.RecordAttributes = helpers.MustUnmarshalJSON[AttributeMap](payload.Record.Attributes)
	encodable.Signatures = payload.Signatures

	return encodable
}

// ToRecordObj converts Record to RecordObj.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (r *ReadableRecord) ToRecordObj() (Record, error) {
	var resourceObj Record

	resourceObj.Id = r.ID
	resourceObj.BondId = r.BondID
	resourceObj.CreateTime = r.CreateTime
	resourceObj.ExpiryTime = r.ExpiryTime
	resourceObj.Deleted = r.Deleted
	resourceObj.Owners = r.Owners
	resourceObj.Attributes = helpers.MustMarshalJSON(r.Attributes)

	return resourceObj, nil
}

// ToReadableRecord converts Record to a serializable object
func (r *Record) ToReadableRecord() ReadableRecord {
	var resourceObj ReadableRecord

	resourceObj.ID = r.Id
	resourceObj.BondID = r.BondId
	resourceObj.CreateTime = r.CreateTime
	resourceObj.ExpiryTime = r.ExpiryTime
	resourceObj.Deleted = r.Deleted
	resourceObj.Owners = r.Owners
	resourceObj.Names = r.Names
	resourceObj.Attributes = helpers.MustUnmarshalJSON[AttributeMap](r.Attributes)

	return resourceObj
}

// CanonicalJSON returns the canonical JSON representation of the record.
func (r *ReadableRecord) CanonicalJSON() []byte {
	bytes, err := canonicaljson.Marshal(r.Attributes)
	if err != nil {
		panic("error marshaling record: " + err.Error())
	}

	return bytes
}

// GetSignBytes generates a record hash to be signed.
func (r *ReadableRecord) GetSignBytes() ([]byte, []byte) {
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
func (r *ReadableRecord) GetCID() (string, error) {
	id, err := helpers.GetCid(r.CanonicalJSON())
	if err != nil {
		return "", err
	}

	return id, nil
}
