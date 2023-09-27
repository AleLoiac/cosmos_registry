package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"errors"
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
	Counter  collections.Map[string, uint64]
	Balances collections.Map[string, uint64]
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
		Counter:      collections.NewMap(sb, example.CounterKey, "counter", collections.StringKey, collections.Uint64Value),
		Balances:     collections.NewMap(sb, example.BalancesKey, "balances", collections.StringKey, collections.Uint64Value),
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

func (k Keeper) MintCoins(ctx context.Context, address string, amount uint64) error {
	balance, err := k.Balances.Get(ctx, address)
	switch {
	case errors.Is(err, collections.ErrNotFound):
		return k.Balances.Set(ctx, address, amount)
	case err != nil:
		return err
	default:
		return k.Balances.Set(ctx, address, balance+amount)
	}

}

func (k Keeper) TransferCoins(ctx context.Context, sender string, receiver string, amount uint64) error {
	senderBalance, err := k.Balances.Get(ctx, sender)
	if err != nil {
		return err
	}
	if senderBalance < amount {
		return errors.New("insufficient balance")
	}
	err = k.Balances.Set(ctx, sender, senderBalance-amount)
	if err != nil {
		return err
	}

	return k.MintCoins(ctx, receiver, amount)
}

func (k Keeper) Burn(ctx context.Context, address string, amount uint64) error {
	balance, err := k.Balances.Get(ctx, address)
	switch {
	case err != nil:
		return err
	case balance < amount:
		return errors.New("insufficient balance")
	default:
		return k.Balances.Set(ctx, address, balance-amount)
	}
}

func walkFunc(key string, value uint64) (stop bool, err error) {
	fmt.Printf("Key: %v, Value: %v\n", key, value)
	return false, nil
}

func (k Keeper) Export(ctx context.Context) error {

	err := k.Balances.Walk(ctx, nil, walkFunc)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	Address string
	Wallet  Coins
}

type Coins struct {
	denomination string
	amount       uint64
}

type Escrow struct {
	creatorAddress string
	blockedToken   Coins
	requestToken   Coins
	claimed        bool
}

type BankKeeper struct {
	escrows map[string]Escrow
}

func CreateEscrow(creator User, blockedToken Coins, requestToken Coins, escrowID string, keeper BankKeeper) error {

	if creator.Wallet.denomination == blockedToken.denomination && creator.Wallet.amount < blockedToken.amount {
		return errors.New("not enough coins")
	}

	escrow := Escrow{
		creatorAddress: creator.Address,
		blockedToken:   blockedToken,
		requestToken:   requestToken,
		claimed:        false,
	}

	keeper.escrows[escrowID] = escrow
	return nil
}

func ClaimEscrow(claimer User, creator User, escrow Escrow) error {
	if claimer.Wallet.denomination == escrow.requestToken.denomination && claimer.Wallet.amount < escrow.requestToken.amount {
		return errors.New("not enough coins")
	}

	claimer.Wallet.amount = claimer.Wallet.amount + escrow.blockedToken.amount
	creator.Wallet.amount = creator.Wallet.amount + escrow.blockedToken.amount
	escrow.claimed = true
	return nil
}
