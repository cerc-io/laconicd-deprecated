//
// Copyright 2020 Wireline, Inc.
//

package utils

import (
	"bytes"
	"errors"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/storage/memstore"

	canonicalJson "github.com/gibson042/canonicaljson-go"
	cbor "github.com/ipfs/go-ipld-cbor"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	mh "github.com/multiformats/go-multihash"
)

func init() {
	multicodec.RegisterEncoder(0x71, dagcbor.Encode)
	multicodec.RegisterDecoder(0x71, dagcbor.Decode)
}

// GenerateHash returns the hash of the canonicalized JSON input.
func GenerateHash(json map[string]interface{}) (string, []byte, error) {
	content, err := canonicalJson.Marshal(json)
	if err != nil {
		return "", nil, err
	}

	//cid, err := CIDFromJSONBytes(content)
	cidString, err := CIDFromJSONBytesUsingIpldPrime(content)
	if err != nil {
		return "", nil, err
	}

	return cidString, content, nil
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

// CIDFromJSONBytesUsingIpldPrime returns CID (dagcbor) for json (as bytes).
func CIDFromJSONBytesUsingIpldPrime(content []byte) (string, error) {
	lsys := cidlink.DefaultLinkSystem()
	var store = memstore.Store{}

	// We want to store the serialized data somewhere.
	// We'll use an in-memory store for this.  (It's a package scoped variable.)
	// You can use any kind of storage system here;
	// or if you need even more control, you could also write a function that conforms to the linking.BlockWriteOpener interface.
	lsys.SetWriteStorage(&store)
	// To create any links, first we need a LinkPrototype.
	// This gathers together any parameters that might be needed when making a link.
	// (For CIDs, the version, the codec, and the multihash type are all parameters we'll need.)
	// Often, you can probably make this a constant for your whole application.
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,    // Usually '1'.
		Codec:    0x71, // 0x71 means "dag-cbor" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhType:   0x13, // 0x20 means "sha2-512" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhLength: 64,   // sha2-512 hash has a 64-byte sum.
	}}

	// And we need some data to link to!  Here's a quick piece of example data:
	n, err := fluent.Build(basicnode.Prototype.Any, func(na fluent.NodeAssembler) {
		na.AssignBytes(content)
	})
	if err != nil {
		return "", err
	}

	// Now: time to apply the LinkSystem, and do the actual store operation!
	lnk, err := lsys.Store(
		linking.LinkContext{}, // The zero value is fine.  Configure it it you want cancellability or other features.
		lp,                    // The LinkPrototype says what codec and hashing to use.
		n,                     // And here's our data.
	)
	if err != nil {
		return "", err
	}
	return lnk.String(), nil
}
