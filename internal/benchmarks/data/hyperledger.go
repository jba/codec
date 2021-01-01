// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

// From https://github.com/hyperledger/fabric
// Specifically, core/ledger/kvledger/tests/sample_data_helper.go.

import (
	"encoding/json"
	"os"
)

var Hyperledger = BenchmarkData{
	"hyperledger",
	func() (interface{}, error) { return hlDecodeJSON("data/ledgerAPIs.json") },
	func() interface{} { return new(submittedData) },
}

type submittedData map[string]*submittedLedgerData

type submittedLedgerData struct {
	Blocks []*BlockAndPvtData
	Txs    []*txAndPvtdata
}

type BlockAndPvtData struct {
	Block          *Block
	PvtData        TxPvtDataMap
	MissingPvtData TxMissingPvtData
}

type TxPvtData struct {
	SeqInBlock uint64
	WriteSet   *TxPvtReadWriteSet
}

type TxPvtReadWriteSet struct {
	DataModel  int32
	NsPvtRwset []*NsPvtReadWriteSet
}

type NsPvtReadWriteSet struct {
	Namespace          string
	CollectionPvtRwset []*CollectionPvtReadWriteSet
}

type CollectionPvtReadWriteSet struct {
	CollectionName string
	Rwset          []byte
}

type TxPvtDataMap map[uint64]*TxPvtData

type Block struct {
	Header   *BlockHeader
	Data     *BlockData
	Metadata *BlockMetadata
}

type BlockHeader struct {
	Number       uint64
	PreviousHash []byte
	DataHash     []byte
}

type BlockData struct {
	Data [][]byte
}

type BlockMetadata struct {
	Metadata [][]byte
}

type txAndPvtdata struct {
	Txid     string
	Envelope *Envelope
	Pvtws    *TxPvtReadWriteSet
}

type Envelope struct {
	Payload   []byte
	Signature []byte
}
type TxMissingPvtData map[uint64][]*MissingPvtData

type MissingPvtData struct {
	Namespace  string
	Collection string
	IsEligible bool
}

func hlDecodeJSON(filename string) (submittedData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	var sd submittedData
	if err := d.Decode(&sd); err != nil {
		return nil, err
	}
	return sd, nil
}
