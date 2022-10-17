//
// Copyright 2020 Wireline, Inc.
//

package utils

import (
	"bytes"
	"encoding/binary"
	"sort"

	set "github.com/deckarep/golang-set"
)

func Int64ToBytes(num int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, num); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func SetToSlice(set set.Set) []string {
	names := []string{}

	for name := range set.Iter() {
		if name, ok := name.(string); ok && name != "" {
			names = append(names, name)
		}
	}

	sort.SliceStable(names, func(i, j int) bool { return names[i] < names[j] })

	return names
}

func SliceToSet(names []string) set.Set {
	set := set.NewThreadUnsafeSet()

	for _, name := range names {
		if name != "" {
			set.Add(name)
		}
	}

	return set
}

func AppendUnique(list []string, element string) []string {
	set := SliceToSet(list)
	set.Add(element)

	return SetToSlice(set)
}
