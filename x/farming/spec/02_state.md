<!-- order: 2 -->

# State

The farming module keeps track of the staking and rewards states.

## Plan Interface

The plan interface exposes methods to read and write standard farming plan information. Note that all of these methods operate on a plan struct confirming to the interface and in order to write the plan to the store, the plan keeper will need to be used.

```go
// PlanI is an interface used to store plan records within state.
type PlanI interface {
    proto.Message

    GetId() uint64
    SetId(uint64) error

    GetType() int32
    SetType(int32) error

    GetFarmingPoolAddress() sdk.AccAddress
    SetFarmingPoolAddress(sdk.AccAddress) error

    GetTerminationAddress() sdk.AccAddress
    SetTerminationAddress(sdk.AccAddress) error

    GetStakingCoinsWeight() sdk.DecCoins
    SetStakingCoinsWeight(sdk.DecCoins) error

    GetStartTime() time.Time
    SetStartTime(time.Time) error

    GetEndTime() time.Time
    SetEndTime(time.Time) error

    GetTerminated() bool
    SetTerminated(bool) error

    GetLastDistributionTime() *time.Time
    SetLastDistributionTime(*time.Time) error

    GetDistributedCoins() sdk.Coins
    SetDistributedCoins(sdk.Coins) error

    String() string
}
```

## Base Plan

A base plan is the simplest and most common plan type, which just stores all requisite fields directly in a struct.

```go
// BasePlan defines a base plan type. It contains all the necessary fields
// for basic farming plan functionality. Any custom farming plan type should extend this
// type for additional functionality (e.g. fixed amount plan, ratio plan).
type BasePlan struct {
    Id                   uint64       // index of the plan
    Type                 PlanType     // type of the plan; public or private
    FarmingPoolAddress   string       // bech32-encoded farming pool address
    TerminationAddress   string       // bech32-encoded termination address
    StakingCoinWeights   sdk.DecCoins // coin weights for the plan
    StartTime            time.Time    // start time of the plan
    EndTime              time.Time    // end time of the plan
    Terminated           bool         // whether the plan has terminated or not
    LastDistributionTime *time.Time   // last time a distribution happened
    DistributedCoins     sdk.Coins    // total coins distributed
}
```

```go
// FixedAmountPlan defines a fixed amount plan that fixed amount of coins are distributed for every epoch day.
type FixedAmountPlan struct {
    *BasePlan

    EpochAmount sdk.Coins // distributing amount for each epoch
}
```

```go
// RatioPlan defines a ratio plan that ratio of total coins in farming pool address is distributed for every epoch day.
type RatioPlan struct {
    *BasePlan

    EpochRatio sdk.Dec // distributing amount by ratio
}
```

## Plan Types

```go
// PlanType enumerates the valid types of a plan.
type PlanType int32

const (
    // PLAN_TYPE_UNSPECIFIED defines the default plan type.
    PlanTypeNil PlanType = 0
    // PLAN_TYPE_PUBLIC defines the public plan type.
    PlanTypePublic PlanType = 1
    // PLAN_TYPE_PRIVATE defines the private plan type.
    PlanTypePrivate PlanType = 2
)
```

The parameters of the Plan state are:

- ModuleName, RouterKey, StoreKey, QuerierRoute: `farming`
- Plan: `0x11 | Id -> ProtocolBuffer(Plan)`
- PlanByFarmerAddrIndex: `0x12 | FarmerAddrLen (1 byte) | FarmerAddr -> Id -> nil` (can be deprecated)
- GlobalPlanIdKey: `[]byte("globalPlanId") -> ProtocolBuffer(uint64)`
  - store latest plan id
- ModuleName, RouterKey, StoreKey, QuerierRoute: `farming`

## Epoch

```go
// meta data for epoch eras
type EpochEra struct {
    EpochEraId  uint64
    EpochPeriod uint64
    StartHeight uint64
    EndHeight   uint64
}
```
- `CurrentEpochEraId` : `EpochEraId` of current epoch era
- `EndHeight` can be 0 if the epoch era is the current one
- new `EpochEra` is created if `EpochPeriod` change proposal is passed
- `EndHeight` is modified when new `EpochEra` is introduced


```go
// meta data for epochs
type Epoch struct {
    EpochEraId  uint64
    EpochId     uint64
    StartHeight uint64
    EndHeight   uint64
}
```
- `CurrentEpochId` : `EpochId` of current epoch
- `EpochId` is always increasing, even if new epoch era starts
- new `Epoch` is created when new epoch starts

## Staking

```go
// Staking stores farmer's staking position status.
type Staking struct {
    Id              uint64
    Farmer          string
    StakedCoins     sdk.Coins
    QueuedCoins     sdk.Coins
    StartEpochId    uint64
}
```

The parameters of the Staking state are:

- GlobalStakingIdKey: `[]byte("globalStakingId") -> ProtocolBuffer(uint64)`

  - store latest staking id

- Staking: `0x21 | Id -> ProtocolBuffer(Staking)`
- StakingByFarmerAddrIndex: `0x22 | FarmerAddrLen (1 byte) | FarmerAddr -> BigEndian(Id)`
- StakingByStakingCoinDenomIdIndex: `0x23 | StakingCoinDenomLen (1 byte) | StakingCoinDenom | BigEndian(Id) -> nil`

## Accumulated Unit Reward(AUR)

```go
// Store AUR for every staking coin and every block height
type AccumulatedUnitReward struct {
    StakingCoinDenom    string
    UntilEpochId        uint64
    AccumulatedReward   sdk.Coins
}
```

- New `AccumulatedUnitReward` struct should be created and managed every start of epoch when new staking coin is introduced from new plan.
- `AccumulatedReward` can be calculated from total block rewards for this staking coin from all existing plan, plus last `AccumulatedReward`


## Examples

An example of `FixedAmountPlan`

```json
{
  "base_plan": {
    "id": 0,
    "type": 0,
    "farmingPoolAddress": "cosmos1...",
    "rewardPoolAddress": "cosmos1...",
    "stakingCoinWeights": [
      {
        "denom": "xxx",
        "amount": "0.200000000000000000"
      },
      {
        "denom": "yyy",
        "amount": "0.300000000000000000"
      },
      {
        "denom": "zzz",
        "amount": "0.500000000000000000"
      }
    ],
    "startTime": "2021-10-01T00:00:00Z",
    "endTime": "2022-04-01T00:00:00Z",
    "terminationAddress": "cosmos1..."
  },
  "epochAmount": {
    "denom": "uatom",
    "amount": "10000000"
  }
}
```

An example of `RatioPlan`

```json
{
  "base_plan": {
    "id": 0,
    "type": 0,
    "farmingPoolAddress": "cosmos1...",
    "rewardPoolAddress": "cosmos1...",
    "stakingCoinWeights": [
      {
        "denom": "xxx",
        "amount": "0.200000000000000000"
      },
      {
        "denom": "yyy",
        "amount": "0.300000000000000000"
      },
      {
        "denom": "zzz",
        "amount": "0.500000000000000000"
      }
    ],
    "startTime": "2021-10-01T00:00:00Z",
    "endTime": "2022-04-01T00:00:00Z",
    "terminationAddress": "cosmos1..."
  },
  "epochRatio": "0.01"
}
```
