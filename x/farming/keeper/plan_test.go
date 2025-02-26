package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming/types"
)

func (suite *KeeperTestSuite) TestGetSetNewPlan() {
	name := ""
	farmingPoolAddr := sdk.AccAddress("farmingPoolAddr")
	terminationAddr := sdk.AccAddress("terminationAddr")

	stakingCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)))
	coinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)

	addrs := app.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(2000000))
	farmerAddr := addrs[0]

	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	basePlan := types.NewBasePlan(1, name, 1, farmingPoolAddr.String(), terminationAddr.String(), coinWeights, startTime, endTime)
	fixedPlan := types.NewFixedAmountPlan(basePlan, sdk.NewCoins(sdk.NewCoin("testFarmCoinDenom", sdk.NewInt(1000000))))
	suite.keeper.SetPlan(suite.ctx, fixedPlan)

	planGet, found := suite.keeper.GetPlan(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(fixedPlan.Id, planGet.GetId())
	suite.Require().Equal(fixedPlan.FarmingPoolAddress, planGet.GetFarmingPoolAddress().String())

	plans := suite.keeper.GetAllPlans(suite.ctx)
	suite.Require().Len(plans, 1)
	suite.Require().Equal(fixedPlan.Id, plans[0].GetId())
	suite.Require().Equal(fixedPlan.FarmingPoolAddress, plans[0].GetFarmingPoolAddress().String())

	_, err := suite.keeper.Stake(suite.ctx, farmerAddr, stakingCoins)
	suite.Require().NoError(err)

	stakings := suite.keeper.GetAllStakings(suite.ctx)
	stakingByFarmer, found := suite.keeper.GetStakingByFarmer(suite.ctx, farmerAddr)
	stakingsByDenom := suite.keeper.GetStakingsByStakingCoinDenom(suite.ctx, sdk.DefaultBondDenom)

	suite.Require().True(found)
	suite.Require().Equal(stakings[0], stakingByFarmer)
	suite.Require().Equal(stakings, stakingsByDenom)
}

func (suite *KeeperTestSuite) TestGlobalPlanId() {
	globalPlanId := suite.keeper.GetGlobalPlanId(suite.ctx)
	suite.Require().Equal(uint64(0), globalPlanId)

	cacheCtx, _ := suite.ctx.CacheContext()
	nextPlanId := suite.keeper.GetNextPlanIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(1), nextPlanId)

	sampleFixedPlan := suite.sampleFixedAmtPlans[0].(*types.FixedAmountPlan)
	poolAcc, err := suite.keeper.GeneratePrivatePlanFarmingPoolAddress(suite.ctx, sampleFixedPlan.Name)
	suite.Require().NoError(err)
	_, err = suite.keeper.CreateFixedAmountPlan(suite.ctx, &types.MsgCreateFixedAmountPlan{
		Name:               sampleFixedPlan.Name,
		Creator:            suite.addrs[0].String(),
		StakingCoinWeights: sampleFixedPlan.GetStakingCoinWeights(),
		StartTime:          sampleFixedPlan.GetStartTime(),
		EndTime:            sampleFixedPlan.GetEndTime(),
		EpochAmount:        sampleFixedPlan.EpochAmount,
	}, poolAcc, suite.addrs[0], types.PlanTypePublic)
	suite.Require().NoError(err)

	globalPlanId = suite.keeper.GetGlobalPlanId(suite.ctx)
	suite.Require().Equal(uint64(1), globalPlanId)

	plans := suite.keeper.GetAllPlans(suite.ctx)
	suite.Require().Len(plans, 1)
	suite.Require().Equal(uint64(len(plans)), globalPlanId)

	cacheCtx, _ = suite.ctx.CacheContext()
	nextPlanId = suite.keeper.GetNextPlanIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(2), nextPlanId)

	sampleRatioPlan := suite.sampleRatioPlans[0].(*types.RatioPlan)
	poolAcc, err = suite.keeper.GeneratePrivatePlanFarmingPoolAddress(suite.ctx, sampleRatioPlan.Name)
	suite.Require().NoError(err)
	_, err = suite.keeper.CreateRatioPlan(suite.ctx, &types.MsgCreateRatioPlan{
		Name:               sampleRatioPlan.Name,
		Creator:            suite.addrs[1].String(),
		StakingCoinWeights: sampleRatioPlan.GetStakingCoinWeights(),
		StartTime:          sampleRatioPlan.GetStartTime(),
		EndTime:            sampleRatioPlan.GetEndTime(),
		EpochRatio:         sampleRatioPlan.EpochRatio,
	}, poolAcc, suite.addrs[1], types.PlanTypePrivate)
	suite.Require().NoError(err)

	globalPlanId = suite.keeper.GetGlobalPlanId(suite.ctx)
	suite.Require().Equal(uint64(2), globalPlanId)

	plans = suite.keeper.GetAllPlans(suite.ctx)
	suite.Require().Len(plans, 2)
	suite.Require().Equal(uint64(len(plans)), globalPlanId)

	cacheCtx, _ = suite.ctx.CacheContext()
	nextPlanId = suite.keeper.GetNextPlanIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(3), nextPlanId)
}
