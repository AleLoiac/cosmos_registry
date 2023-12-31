package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmosregistry/example"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Schema   collections.Schema
	Params   collections.Item[example.Params]
	TweetsID collections.Sequence
	Tweets   collections.Map[uint64, example.Tweet]
	LikedBy  collections.Map[uint64, string]
	Counter  collections.Map[string, uint64]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, example.ParamsKey, "params", codec.CollValue[example.Params](cdc)),
		TweetsID:     collections.NewSequence(sb, example.TweetsIDKey, "tweets_id"),
		Tweets:       collections.NewMap(sb, example.TweetsKey, "tweets", collections.Uint64Key, codec.CollValue[example.Tweet](cdc)),
		LikedBy:      collections.NewMap(sb, example.LikedByKey, "liked_by", collections.Uint64Key, collections.StringValue),
		Counter:      collections.NewMap(sb, example.CounterKey, "counter", collections.StringKey, collections.Uint64Value),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
