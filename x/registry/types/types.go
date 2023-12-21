package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cerc-io/laconicd/x/registry/helpers"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	canonicalJson "github.com/gibson042/canonicaljson-go"
	"github.com/gogo/protobuf/proto"
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
	payload := Payload{
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
	recordType, ok := recordPayLoad["type"]
	if !ok {
		return &codectypes.Any{}, fmt.Errorf("cannot get type from payload")
	}
	bz := helpers.MarshalMapToJSONBytes(recordPayLoad)

	switch recordType.(string) {
	case "ServiceProviderRegistration":
		{
			var attributes ServiceProviderRegistration
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "WebsiteRegistrationRecord":
		{
			var attributes WebsiteRegistrationRecord
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "ApplicationRecord":
		{
			var attributes ApplicationRecord
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "ApplicationDeploymentRequest":
		{
			var attributes ApplicationDeploymentRequest
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "ApplicationDeploymentRecord":
		{
			var attributes ApplicationDeploymentRecord
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "ApplicationDeploymentRemovalRequest":
		{
			var attributes ApplicationDeploymentRemovalRequest
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "ApplicationDeploymentRemovalRecord":
		{
			var attributes ApplicationDeploymentRemovalRecord
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "ApplicationArtifact":
		{
			var attributes ApplicationArtifact
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "DnsRecord":
		{
			var attributes DnsRecord
			err := json.Unmarshal(bz, &attributes)
			if err != nil {
				return &codectypes.Any{}, err
			}
			return codectypes.NewAnyWithValue(&attributes)
		}
	case "GeneralRecord":
		{
			var attributes GeneralRecord
			err := json.Unmarshal(bz, &attributes)
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
	bz, err := GetJSONBytesFromAny(*payload.Record.Attributes)
	if err != nil {
		panic(err)
	}

	payloadType.Record = helpers.UnMarshalMapFromJSONBytes(bz)
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

	bz, err := GetJSONBytesFromAny(*r.Attributes)
	if err != nil {
		panic(err)
	}
	resourceObj.Attributes = helpers.UnMarshalMapFromJSONBytes(bz)

	return resourceObj
}

func GetJSONBytesFromAny(any codectypes.Any) ([]byte, error) {
	var bz []byte
	s := strings.Split(any.TypeUrl, ".")
	switch s[len(s)-1] {
	case "ServiceProviderRegistration":
		{
			var attributes ServiceProviderRegistration
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "WebsiteRegistrationRecord":
		{
			var attributes WebsiteRegistrationRecord
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "ApplicationRecord":
		{
			var attributes ApplicationRecord
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "ApplicationDeploymentRequest":
		{
			var attributes ApplicationDeploymentRequest
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "ApplicationDeploymentRecord":
		{
			var attributes ApplicationDeploymentRecord
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "ApplicationDeploymentRemovalRequest":
		{
			var attributes ApplicationDeploymentRemovalRequest
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "ApplicationDeploymentRemovalRecord":
		{
			var attributes ApplicationDeploymentRemovalRecord
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "ApplicationArtifact":
		{
			var attributes ApplicationArtifact
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "DnsRecord":
		{
			var attributes DnsRecord
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	case "GeneralRecord":
		{
			var attributes GeneralRecord
			err := proto.Unmarshal(any.Value, &attributes)
			if err != nil {
				panic("Proto unmarshal error")
			}

			bz, err = json.Marshal(attributes)
			if err != nil {
				panic("JSON marshal error")
			}
		}
	default:
		return nil, fmt.Errorf("unsupported type %s", s[len(s)-1])
	}

	return bz, nil
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
