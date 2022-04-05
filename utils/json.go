//
// Copyright 2020 Wireline, Inc.
//

package utils

import (
	"bytes"
	"errors"

	canonicalJson "github.com/gibson042/canonicaljson-go"
	cbor "github.com/ipfs/go-ipld-cbor"
	mh "github.com/multiformats/go-multihash"
)

// GenerateHash returns the hash of the canonicalized JSON input.
func GenerateHash(json map[string]interface{}) (string, []byte, error) {
	content, err := canonicalJson.Marshal(json)
	if err != nil {
		return "", nil, err
	}

	cid, err := CIDFromJSONBytes(content)
	if err != nil {
		return "", nil, err
	}

	return cid, content, nil
}

// CIDFromJSONBytes returns CID (cbor) for json (as bytes).
func CIDFromJSONBytes(content []byte) (string, error) {
	cid, err := cbor.FromJSON(bytes.NewReader(content), mh.SHA2_256, -1)
	if err != nil {
		return "", err
	}

	return cid.String(), nil
}

// GetAttributeAsString returns a map attribute as string, if possible.
func GetAttributeAsString(obj map[string]interface{}, attr string) (string, error) {
	if value, ok := obj[attr]; ok {
		if valueStr, ok := value.(string); ok {
			return valueStr, nil
		}

		return "", errors.New("attribute not of string type")
	}

	return "", errors.New("attribute not found")
}
