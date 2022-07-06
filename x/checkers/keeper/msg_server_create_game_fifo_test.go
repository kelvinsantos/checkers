package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCreate3GamesHasSavedFifo(t *testing.T) {
	msgSrvr, keeper, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
	})

	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	nextGame2, found2 := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, found2)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  3,
		FifoHead: "1",
		FifoTail: "2",
	}, nextGame2)
	game1, found1 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, found1)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(0),
		BeforeId:  "-1",
		AfterId:   "2",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game1)
	game2, found2 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "2")
	require.True(t, found2)
	require.EqualValues(t, types.StoredGame{
		Creator:   bob,
		Index:     "2",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       carol,
		Black:     alice,
		MoveCount: uint64(0),
		BeforeId:  "1",
		AfterId:   "-1",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game2)

	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Red:     alice,
		Black:   bob,
	})
	nextGame3, found3 := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, found3)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  4,
		FifoHead: "1",
		FifoTail: "3",
	}, nextGame3)
	game1, found1 = keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, found1)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(0),
		BeforeId:  "-1",
		AfterId:   "2",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game1)
	game2, found2 = keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "2")
	require.True(t, found2)
	require.EqualValues(t, types.StoredGame{
		Creator:   bob,
		Index:     "2",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       carol,
		Black:     alice,
		MoveCount: uint64(0),
		BeforeId:  "1",
		AfterId:   "3",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game2)
	game3, found3 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "3")
	require.True(t, found3)
	require.EqualValues(t, types.StoredGame{
		Creator:   carol,
		Index:     "3",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       alice,
		Black:     bob,
		MoveCount: uint64(0),
		BeforeId:  "2",
		AfterId:   "-1",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game3)
}

func TestPlayMove2Games2MovesHasSavedFifo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: alice,
		IdValue: "2",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	nextGame1, found1 := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, found1)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  3,
		FifoHead: "1",
		FifoTail: "2",
	}, nextGame1)
	game1, found1 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, found1)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "r",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(1),
		BeforeId:  "-1",
		AfterId:   "2",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game1)
	game2, found2 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "2")
	require.True(t, found2)
	require.EqualValues(t, types.StoredGame{
		Creator:   bob,
		Index:     "2",
		Game:      "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "r",
		Red:       carol,
		Black:     alice,
		MoveCount: uint64(1),
		BeforeId:  "1",
		AfterId:   "-1",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game2)
}

func TestRejectMiddleGameHasSavedFifo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Red:     alice,
		Black:   bob,
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator: carol,
		IdValue: "2",
	})
	nextGame, found := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  4,
		FifoHead: "1",
		FifoTail: "3",
	}, nextGame)
	game1, found1 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, found1)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(0),
		BeforeId:  "-1",
		AfterId:   "3",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game1)
	game3, found3 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "3")
	require.True(t, found3)
	require.EqualValues(t, types.StoredGame{
		Creator:   carol,
		Index:     "3",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       alice,
		Black:     bob,
		MoveCount: uint64(0),
		BeforeId:  "1",
		AfterId:   "-1",
		Deadline:  types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:    "*",
	}, game3)
}
