package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateFixedAmountPlan)(nil)
	_ sdk.Msg = (*MsgCreateRatioPlan)(nil)
	_ sdk.Msg = (*MsgStake)(nil)
	_ sdk.Msg = (*MsgUnstake)(nil)
	_ sdk.Msg = (*MsgHarvest)(nil)
)

// Message types for the farming module
const (
	TypeMsgCreateFixedAmountPlan = "create_fixed_amount_plan"
	TypeMsgCreateRatioPlan       = "create_ratio_plan"
	TypeMsgStake                 = "stake"
	TypeMsgUnstake               = "unstake"
	TypeMsgHarvest               = "harvest"
)

// NewMsgCreateFixedAmountPlan creates a new MsgCreateFixedAmountPlan.
func NewMsgCreateFixedAmountPlan(
	name string,
	creatorAcc sdk.AccAddress,
	stakingCoinWeights sdk.DecCoins,
	startTime time.Time,
	endTime time.Time,
	epochAmount sdk.Coins,
) *MsgCreateFixedAmountPlan {
	return &MsgCreateFixedAmountPlan{
		Name:               name,
		Creator:            creatorAcc.String(),
		StakingCoinWeights: stakingCoinWeights,
		StartTime:          startTime,
		EndTime:            endTime,
		EpochAmount:        epochAmount,
	}
}

func (msg MsgCreateFixedAmountPlan) Route() string { return RouterKey }

func (msg MsgCreateFixedAmountPlan) Type() string { return TypeMsgCreateFixedAmountPlan }

func (msg MsgCreateFixedAmountPlan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address %q: %v", msg.Creator, err)
	}
	if !msg.EndTime.After(msg.StartTime) {
		return sdkerrors.Wrapf(ErrInvalidPlanEndTime, "end time %s must be greater than start time %s", msg.EndTime.Format(time.RFC3339), msg.StartTime.Format(time.RFC3339))
	}
	if msg.StakingCoinWeights.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking coin weights must not be empty")
	}
	if err := msg.StakingCoinWeights.Validate(); err != nil {
		return err
	}
	if ok := ValidateStakingCoinTotalWeights(msg.StakingCoinWeights); !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "total weight must be 1")
	}
	if msg.EpochAmount.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "epoch amount must not be empty")
	}
	if err := msg.EpochAmount.Validate(); err != nil {
		return err
	}
	return nil
}

func (msg MsgCreateFixedAmountPlan) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateFixedAmountPlan) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateFixedAmountPlan) GetCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCreateRatioPlan creates a new MsgCreateRatioPlan.
func NewMsgCreateRatioPlan(
	name string,
	CreatorAcc sdk.AccAddress,
	stakingCoinWeights sdk.DecCoins,
	startTime time.Time,
	endTime time.Time,
	epochRatio sdk.Dec,
) *MsgCreateRatioPlan {
	return &MsgCreateRatioPlan{
		Name:               name,
		Creator:            CreatorAcc.String(),
		StakingCoinWeights: stakingCoinWeights,
		StartTime:          startTime,
		EndTime:            endTime,
		EpochRatio:         epochRatio,
	}
}

func (msg MsgCreateRatioPlan) Route() string { return RouterKey }

func (msg MsgCreateRatioPlan) Type() string { return TypeMsgCreateRatioPlan }

func (msg MsgCreateRatioPlan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address %q: %v", msg.Creator, err)
	}
	if !msg.EndTime.After(msg.StartTime) {
		return sdkerrors.Wrapf(ErrInvalidPlanEndTime, "end time %s must be greater than start time %s", msg.EndTime.Format(time.RFC3339), msg.StartTime.Format(time.RFC3339))
	}
	if msg.StakingCoinWeights.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking coin weights must not be empty")
	}
	if err := msg.StakingCoinWeights.Validate(); err != nil {
		return err
	}
	if ok := ValidateStakingCoinTotalWeights(msg.StakingCoinWeights); !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "total weight must be 1")
	}
	if !msg.EpochRatio.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid epoch ratio")
	}
	if msg.EpochRatio.GT(sdk.NewDec(1)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid epoch ratio")
	}
	return nil
}

func (msg MsgCreateRatioPlan) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateRatioPlan) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateRatioPlan) GetCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgStake creates a new MsgStake.
func NewMsgStake(
	farmer sdk.AccAddress,
	stakingCoins sdk.Coins,
) *MsgStake {
	return &MsgStake{
		Farmer:       farmer.String(),
		StakingCoins: stakingCoins,
	}
}

func (msg MsgStake) Route() string { return RouterKey }

func (msg MsgStake) Type() string { return TypeMsgStake }

func (msg MsgStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid farmer address %q: %v", msg.Farmer, err)
	}
	if ok := msg.StakingCoins.IsZero(); ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking coins must not be zero")
	}
	if err := msg.StakingCoins.Validate(); err != nil {
		return err
	}
	return nil
}

func (msg MsgStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgStake) GetFarmer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgUnstake creates a new MsgUnstake.
func NewMsgUnstake(
	farmer sdk.AccAddress,
	unstakingCoins sdk.Coins,
) *MsgUnstake {
	return &MsgUnstake{
		Farmer:         farmer.String(),
		UnstakingCoins: unstakingCoins,
	}
}

func (msg MsgUnstake) Route() string { return RouterKey }

func (msg MsgUnstake) Type() string { return TypeMsgUnstake }

func (msg MsgUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid farmer address %q: %v", msg.Farmer, err)
	}
	if ok := msg.UnstakingCoins.IsZero(); ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unstaking coins must not be zero")
	}
	if err := msg.UnstakingCoins.Validate(); err != nil {
		return err
	}
	return nil
}

func (msg MsgUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgUnstake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgUnstake) GetFarmer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgHarvest creates a new MsgHarvest.
func NewMsgHarvest(
	farmer sdk.AccAddress,
	stakingCoinDenoms []string,
) *MsgHarvest {
	return &MsgHarvest{
		Farmer:            farmer.String(),
		StakingCoinDenoms: stakingCoinDenoms,
	}
}

func (msg MsgHarvest) Route() string { return RouterKey }

func (msg MsgHarvest) Type() string { return TypeMsgHarvest }

func (msg MsgHarvest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid farmer address %q: %v", msg.Farmer, err)
	}
	if len(msg.StakingCoinDenoms) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking coin denoms must be provided at least one")
	}
	for _, denom := range msg.StakingCoinDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return err
		}
	}
	return nil
}

func (msg MsgHarvest) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgHarvest) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgHarvest) GetFarmer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}
