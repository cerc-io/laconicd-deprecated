//
// Copyright 2020 Wireline, Inc.
//

package utils

import (
	"bytes"
	"errors"

	canonicalJson "github.com/gibson042/canonicaljson-go"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/multicodec"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	mh "github.com/multiformats/go-multihash"
)

var store = memstore.Store{}

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
// This is combination of samples for unmarshalling and linking
// see: https://pkg.go.dev/github.com/ipld/go-ipld-prime
func CIDFromJSONBytesUsingIpldPrime(content []byte) (string, error) {
	np := basicnode.Prototype.Any                       // Pick a stle for the in-memory data.
	nb := np.NewBuilder()                               // Create a builder.
	err := dagjson.Decode(nb, bytes.NewReader(content)) // Hand the builder to decoding -- decoding will fill it in!
	if err != nil {
		return "", err
	}
	n := nb.Build() // Call 'Build' to get the resulting Node.  (It's immutable!)

	lsys := cidlink.DefaultLinkSystem()

	// We want to store the serialized data somewhere.
	// We'll use an in-memory store for this.  (It's a package scoped variable.)
	// You can use any kind of storage system here;
	// or if you need even more control, you could also write a function that conforms to the linking.BlockWriteOpener interface.
	lsys.SetWriteStorage(&store)
	// To create any links, first we need a LinkPrototype.
	// This gathers together any parameters that might be needed when making a link.
	// (For CIDs, the version, the codec, and the multihash type are all parameters we'll need.)
	// Often, you can probably make this a constant for your whole application.
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{ //nolint:golint
		Version:  1,    // Usually '1'.
		Codec:    0x71, // 0x71 means "dag-cbor" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhType:   0x12, // 0x12 means "sha2-256" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhLength: 32,   // sha2-256 hash has a 32-byte sum.
	}}

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
