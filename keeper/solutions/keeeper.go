package solutions

import (
	"context"
	"cosmossdk.io/collections"
	"fmt"
)

var moduleName = "escrow"

type Coin struct {
	Denom  string
	Amount uint64
}

type Escrow struct {
	LockedCoins Coin
	WantCoin    Coin
}

type BankKeeper interface {
	SendCoins(ctx context.Context, from, to string, amt Coin) error
}

type Keeper struct {
	bk      BankKeeper
	Escrows collections.Map[string, Escrow]
}

func (k Keeper) CreateEscrow(ctx context.Context, creator string, lockedAmount, wantAmount Coin) error {
	exists, err := k.Escrows.Has(ctx, creator)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("can't have multiple escrows")
	}

	escrow := Escrow{
		LockedCoins: lockedAmount,
		WantCoin:    wantAmount,
	}
	err = k.Escrows.Set(ctx, creator, escrow)
	if err != nil {
		return err
	}

	err = k.bk.SendCoins(ctx, creator, moduleName, lockedAmount)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) ClaimEscrow(ctx context.Context, claimer, locker string) error {
	escrow, err := k.Escrows.Get(ctx, locker)
	if err != nil {
		return err
	}

	err = k.bk.SendCoins(ctx, claimer, locker, escrow.WantCoin)
	if err != nil {
		return err
	}

	err = k.bk.SendCoins(ctx, moduleName, claimer, escrow.LockedCoins)
	if err != nil {
		return err
	}

	return k.Escrows.Remove(ctx, locker)
}
