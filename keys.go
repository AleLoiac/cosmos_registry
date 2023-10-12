package example

import "cosmossdk.io/collections"

const ModuleName = "example"

var (
	ParamsKey   = collections.NewPrefix(0)
	CounterKey  = collections.NewPrefix(1)
	TweetsIDKey = collections.NewPrefix(2)
	TweetsKey   = collections.NewPrefix(3)
	LikedByKey  = collections.NewPrefix(4)
)
