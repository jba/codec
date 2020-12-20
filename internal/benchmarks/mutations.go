package main

import "time"

// TODO: use this
// Adapted from https://github.com/luci/luci-go/blob/f41ecd629b2a49d7affcf2729726a293fcf70d32/tumble/model_mutation.go#L119.

type Mutation struct {
	Kind   string
	ID     string
	Parent *Key

	ExpandedShard int64
	ProcessAfter  time.Time
	TargetRoot    *Key

	Version string
	Type    string
	Data    []byte
}

type Key struct {
	Kind      string
	ID        int64
	Name      string
	Parent    *Key
	Namespace string
}
