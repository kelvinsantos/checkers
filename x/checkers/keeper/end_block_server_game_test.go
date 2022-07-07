package keeper_test

import (
	"time"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestForfeitUnplayed() {
	suite.setupSuiteWithOneGameForPlayMove()
	goCtx := sdk.WrapSDKContext(suite.ctx)

	keeper := suite.app.CheckersKeeper
	game1, found := keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().True(found)
	game1.Deadline = types.FormatDeadline(suite.ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(suite.ctx, game1)
	keeper.ForfeitExpiredGames(goCtx)

	_, found = keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().False(found)

	nextGame, found := keeper.GetNextGame(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.NextGame{
		Creator:  "",
		IdValue:  2,
		FifoHead: "-1",
		FifoTail: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(suite.ctx.EventManager().ABCIEvents())
	suite.Require().Len(events, 1)

	forfeitEvent := events[0]
	suite.Require().Equal(forfeitEvent.Type, "message")
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, forfeitEvent.Attributes[createEventCount:])

	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
}

func (suite *IntegrationTestSuite) TestForfeitPlayedOnce() {
	suite.setupSuiteWithOneGameForPlayMove()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	keeper := suite.app.CheckersKeeper
	game1, found := keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().True(found)
	game1.Deadline = types.FormatDeadline(suite.ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(suite.ctx, game1)
	keeper.ForfeitExpiredGames(goCtx)
	suite.RequireBankBalance(balCarol, carol) // Refunded

	_, found = keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().False(found)

	nextGame, found := keeper.GetNextGame(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.NextGame{
		Creator:  "",
		IdValue:  2,
		FifoHead: "-1",
		FifoTail: "-1",
	}, nextGame)

	events := sdk.StringifyEvents(suite.ctx.EventManager().ABCIEvents())
	suite.Require().Len(events, 2)
	forfeitEvent := events[0]
	suite.Require().Equal(forfeitEvent.Type, "message")
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "sender", Value: checkersModuleAddress},
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, forfeitEvent.Attributes[createEventCount+playEventCountFirst:])

	transferEvent := events[1]
	suite.Require().Equal(transferEvent.Type, "transfer")
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "recipient", Value: carol},
		{Key: "sender", Value: checkersModuleAddress},
		{Key: "amount", Value: "11stake"},
	}, transferEvent.Attributes[transferEventCount:])

	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
}

func (suite *IntegrationTestSuite) TestForfeitPlayedTwice() {
	suite.setupSuiteWithOneGameForPlayMove()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.RequireBankBalance(balCarol, carol)
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	suite.RequireBankBalance(balBob, bob)
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   0,
		FromY:   5,
		ToX:     1,
		ToY:     4,
	})
	keeper := suite.app.CheckersKeeper
	game1, found := keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().True(found)
	oldDeadline := types.FormatDeadline(suite.ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(suite.ctx, game1)
	keeper.ForfeitExpiredGames(goCtx)
	suite.RequireBankBalance(balBob+11, bob)     // Won wager
	suite.RequireBankBalance(balCarol-11, carol) // Lost wager

	game1, found = keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(2),
		BeforeId:  "-1",
		AfterId:   "-1",
		Deadline:  oldDeadline,
		Winner:    "r",
		Wager:     11,
	}, game1)

	nextGame, found := keeper.GetNextGame(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.NextGame{
		Creator:  "",
		IdValue:  2,
		FifoHead: "-1",
		FifoTail: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(suite.ctx.EventManager().ABCIEvents())
	suite.Require().Len(events, 2)

	forfeitEvent := events[0]
	suite.Require().Equal(forfeitEvent.Type, "message")
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "sender", Value: checkersModuleAddress},
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "r"},
	}, forfeitEvent.Attributes[createEventCount+2*playEventCountFirst:])

	transferEvent := events[1]
	suite.Require().Equal(transferEvent.Type, "transfer")
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "recipient", Value: bob},
		{Key: "sender", Value: checkersModuleAddress},
		{Key: "amount", Value: "22stake"},
	}, transferEvent.Attributes[2*transferEventCount:])
}

func (suite *IntegrationTestSuite) TestForfeitPlayedTwiceForeignToken() {
	suite.setupSuiteWithBalances()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateGame(goCtx, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
		Wager:   1,
		Token:   foreignToken,
	})
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	suite.RequireBankBalance(balBob, bob)
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   0,
		FromY:   5,
		ToX:     1,
		ToY:     4,
	})
	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
	suite.RequireBankBalanceIn(balTokenAlice, alice, foreignToken)
	suite.RequireBankBalanceIn(balTokenBob-1, bob, foreignToken)
	suite.RequireBankBalanceIn(balTokenCarol-1, carol, foreignToken)
	suite.RequireBankBalanceIn(2, checkersModuleAddress, foreignToken)
	keeper := suite.app.CheckersKeeper
	game1, found := keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().True(found)
	oldDeadline := types.FormatDeadline(suite.ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(suite.ctx, game1)
	keeper.ForfeitExpiredGames(goCtx)
	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
	suite.RequireBankBalanceIn(balTokenAlice, alice, foreignToken)
	suite.RequireBankBalanceIn(balTokenBob+1, bob, foreignToken)     // Won wager
	suite.RequireBankBalanceIn(balTokenCarol-1, carol, foreignToken) // Lost wager
	suite.RequireBankBalanceIn(0, checkersModuleAddress, foreignToken)
}
