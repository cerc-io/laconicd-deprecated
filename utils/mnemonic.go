//
// Copyright 2020 Wireline, Inc.
//

package utils

import "github.com/cosmos/go-bip39"

const (
	mnemonicEntropySize = 256
)

func GenerateMnemonic() (string, error) {
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
