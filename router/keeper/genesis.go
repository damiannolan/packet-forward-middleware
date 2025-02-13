package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/packet-forward-middleware/v6/router/types"
)

// InitGenesis
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	k.SetParams(ctx, state.Params)

	// Initialize store refund path for forwarded packets in genesis state that have not yet been acked.
	store := ctx.KVStore(k.storeKey)
	for key, value := range state.InFlightPackets {
		bz := k.cdc.MustMarshal(&value)
		store.Set([]byte(key), bz)
	}
}

// ExportGenesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	store := ctx.KVStore(k.storeKey)

	inFlightPackets := make(map[string]types.InFlightPacket)

	itr := store.Iterator(nil, nil)
	for ; itr.Valid(); itr.Next() {
		var inFlightPacket types.InFlightPacket
		k.cdc.MustUnmarshal(itr.Value(), &inFlightPacket)
		inFlightPackets[string(itr.Key())] = inFlightPacket
	}
	return &types.GenesisState{Params: k.GetParams(ctx), InFlightPackets: inFlightPackets}
}
