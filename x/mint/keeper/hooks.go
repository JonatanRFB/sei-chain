package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochTypes "github.com/sei-protocol/sei-chain/x/epoch/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epoch epochTypes.Epoch) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epoch epochTypes.Epoch) {
	latestMinter := k.GetOrUpdateLatestMinter(ctx, epoch)
	coinsToMint := latestMinter.GetReleaseAmountToday()

	if coinsToMint.IsZero() {
		k.Logger(ctx).Debug("No coins to mint", "minter", latestMinter)
		return
	}

	// mint coins, update supply
	if err := k.MintCoins(ctx, coinsToMint); err != nil {
		panic(err)
	}
	// send the minted coins to the fee collector account
	if err := k.AddCollectedFees(ctx, coinsToMint); err != nil {
		panic(err)
	}

	// Released Succssfully, decrement the remaining amount by the daily release amount and update minter
	latestMinter.RecordSuccessfulMint(ctx, epoch)
	k.SetMinter(ctx, latestMinter)
}

type Hooks struct {
	k Keeper
}

var _ epochTypes.EpochHooks = Hooks{}

// Return the wrapper struct.
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks.
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epoch epochTypes.Epoch) {
	h.k.BeforeEpochStart(ctx, epoch)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epoch epochTypes.Epoch) {
	h.k.AfterEpochEnd(ctx, epoch)
}
