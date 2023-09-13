package example

import "cosmossdk.io/collections"

const ModuleName = "example"

var (
	ParamsKey   = collections.NewPrefix(0)
	CounterKey  = collections.NewPrefix(1)
	BalancesKey = collections.NewPrefix(2)
)
