// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestEventId(t *testing.T) {
	var table = []struct {
		definition   string
		expectations map[string]common.Hash
	}{
		{
			definition: `[
			{ "type" : "event", "name" : "balance", "inputs": [{ "name" : "in", "type": "uint256" }] },
			{ "type" : "event", "name" : "check", "inputs": [{ "name" : "t", "type": "address" }, { "name": "b", "type": "uint256" }] }
			]`,
			expectations: map[string]common.Hash{
				"balance": crypto.Keccak256Hash([]byte("balance(uint256)")),
				"check":   crypto.Keccak256Hash([]byte("check(address,uint256)")),
			},
		},
	}

	for _, test := range table {
		abi, err := JSON(strings.NewReader(test.definition))
		if err != nil {
			t.Fatal(err)
		}

		for name, event := range abi.Events {
			if event.Id() != test.expectations[name] {
				t.Errorf("expected id to be %x, got %x", test.expectations[name], event.Id())
			}
		}
	}
}

func TestEventUnpacking(t *testing.T) {
	for i, test := range []struct {
		definition     string      //abi definition
		data           []byte      //log data gotten from the event
		expectedOutput interface{} // the expected output
		err            string      // empty or error if expected
	}{
		{
			`[{"anonymous":false,"inputs":[{"indexed":false,"name":"a","type":"uint256"},{"indexed":true,"name":"b","type":"uint256"},{"indexed":false,"name":"c","type":"uint256"}],"name":"A","type":"event"}]`,
			append(pad([]byte{1}, 32, true), pad([]byte{3}, 32, true)...),
			[]interface{}{*big.Int{1}, *big.Int{3}},
			"",
		},
		{
			`[{"anonymous":false,"inputs":[{"indexed":false,"name":"a","type":"int256"},{"indexed":true,"name":"b","type":"int256"},{"indexed":false,"name":"c","type":"int256"}],"name":"B","type":"event"}]`,
			append(pad([]byte{1}, 32, true), pad([]byte{3}, 32, true)...),
			[]interface{}{*big.Int{1}, *big.Int{3}},
			"",
		},
		{
			`[{"anonymous":false,"inputs":[{"indexed":false,"name":"a","type":"bool"},{"indexed":false,"name":"b","type":"int256"},{"indexed":false,"name":"c","type":"uint256"}],"name":"C","type":"event"}]`,
			append(pad([]byte{0}, 32, true), append(pad([]byte{1}, 32, true), pad([]byte{3}, 32, true)...)...),
			[]interface{}{false, *big.Int{1}, *big.Int{3}},
			"",
		},
		{
			`[{"anonymous":false,"inputs":[{"indexed":false,"name":"a","type":"string"},{"indexed":false,"name":"b","type":"string"},{"indexed":false,"name":"c","type":"string"}],"name":"D","type":"event"}]`,
		},
		{
			`[{"anonymous":false,"inputs":[{"indexed":true,"name":"a","type":"string"},{"indexed":true,"name":"b","type":"string"},{"indexed":true,"name":"c","type":"string"}],"name":"E","type":"event"}]`,
			[]byte(nil),
			[]interface{}{},
			"",
		},
	} {
		abi, err := JSON(strings.NewReader(test.definition))
		if err != nil {
			t.Fatal(err)
		}

	}
}
