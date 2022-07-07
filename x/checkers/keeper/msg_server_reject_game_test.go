package keeper_test

import (
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) setupSuiteWithOneGameForRejectGame() {
	suite.setupSuiteWithBalances()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateGame(goCtx, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
		Wager:   11,
	})
}

func (suite *IntegrationTestSuite) TestRejectGameByRedOneMoveRefundedForeignToken() {
	suite.setupSuiteWithOneGameForRejectGame()
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
		IdValue: "2",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
	suite.RequireBankBalanceIn(balTokenAlice, alice, foreignToken)
	suite.RequireBankBalanceIn(balTokenBob, bob, foreignToken)
	suite.RequireBankBalanceIn(balTokenCarol-1, carol, foreignToken)
	suite.RequireBankBalanceIn(1, checkersModuleAddress, foreignToken)
	suite.msgServer.RejectGame(goCtx, &types.MsgRejectGame{
		Creator: bob,
		IdValue: "2",
	})
	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
	suite.RequireBankBalanceIn(balTokenAlice, alice, foreignToken)
	suite.RequireBankBalanceIn(balTokenBob, bob, foreignToken)
	suite.RequireBankBalanceIn(balTokenCarol, carol, foreignToken)
	suite.RequireBankBalanceIn(0, checkersModuleAddress, foreignToken)
}

func (suite *IntegrationTestSuite) TestRejectGameByRedOneMoveRemovedGame() {
	suite.setupSuiteWithOneGameForRejectGame()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	suite.msgServer.RejectGame(goCtx, &types.MsgRejectGame{
		Creator: bob,
		IdValue: "1",
	})
	keeper := suite.app.CheckersKeeper
	nextGame, found := keeper.GetNextGame(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.NextGame{
		Creator:  "",
		IdValue:  2,
		FifoHead: "-1",
		FifoTail: "-1",
	}, nextGame)
	_, found = keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().False(found)
}

func (suite *IntegrationTestSuite) TestRejectGameByRedOneMoveEmitted() {
	suite.setupSuiteWithOneGameForRejectGame()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	suite.msgServer.RejectGame(goCtx, &types.MsgRejectGame{
		Creator: bob,
		IdValue: "1",
	})
	events := sdk.StringifyEvents(suite.ctx.EventManager().ABCIEvents())
	suite.Require().Len(events, 2)

	rejectEvent := events[0]
	suite.Require().Equal(rejectEvent.Type, "message")
	rejectAttributesDiscardCount := createEventCount + playEventCountFirst
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "sender", Value: checkersModuleAddress},
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameRejected"},
		{Key: "Creator", Value: bob},
		{Key: "IdValue", Value: "1"},
	}, rejectEvent.Attributes[rejectAttributesDiscardCount:])

	transferEvent := events[1]
	suite.Require().Equal(transferEvent.Type, "transfer")
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "recipient", Value: carol},
		{Key: "sender", Value: checkersModuleAddress},
		{Key: "amount", Value: "11stake"},
	}, transferEvent.Attributes[transferEventCount:])
}

func (suite *IntegrationTestSuite) TestRejectGameByBlackConsumedGas() {
	suite.setupSuiteWithOneGameForRejectGame()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	gasBefore := suite.ctx.GasMeter().GasConsumed()
	suite.msgServer.RejectGame(goCtx, &types.MsgRejectGame{
		Creator: carol,
		IdValue: "1",
	})
	gasAfter := suite.ctx.GasMeter().GasConsumed()
	suite.Require().Equal(uint64(6_116+0), gasAfter-gasBefore)
}
