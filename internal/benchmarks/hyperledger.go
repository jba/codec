package main

// From https://github.com/hyperledger/fabric
// Specifically, core/ledger/kvledger/tests/sample_data_helper.go.

import (
	"encoding/json"
	"os"
)

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
	NsPvtRwset []*NsPvtReadWriteSet `protobuf:"bytes,2,rep,name=ns_pvt_rwset,json=nsPvtRwset,proto3" json:"ns_pvt_rwset,omitempty"`
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
