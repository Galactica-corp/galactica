package distribution

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ErrNotRunInEvm is raised when a function is not called inside the EVM.
	ErrNotRunInEvm = "not run in EVM"
	// ErrDifferentOrigin is raised when an approval is set but the origin address is not the same as the spender.
	ErrDifferentOrigin = "tx origin address %s does not match the delegator address %s"
	// ErrInvalidABI is raised when the ABI cannot be parsed.
	ErrInvalidABI = "invalid ABI: %w"
	// ErrInvalidAmount is raised when the amount cannot be cast to a big.Int.
	ErrInvalidAmount = "invalid amount: %v"
	// ErrInvalidDelegator is raised when the delegator address is not valid.
	ErrInvalidDelegator = "invalid delegator address: %s"
	// ErrInvalidDenom is raised when the denom is not valid.
	ErrInvalidDenom = "invalid denom: %s"
	// ErrInvalidMsgType is raised when the transaction type is not valid for the given precompile.
	ErrInvalidMsgType = "invalid %s transaction type: %s"
	// ErrInvalidNumberOfArgs is raised when the number of arguments is not what is expected.
	ErrInvalidNumberOfArgs = "invalid number of arguments; expected %d; got: %d"
	// ErrUnknownMethod is raised when the method is not known.
	ErrUnknownMethod = "unknown method: %s"
	// ErrIntegerOverflow is raised when an integer overflow occurs.
	ErrIntegerOverflow = "integer overflow"
	// ErrNegativeAmount is raised when an amount is negative.
	ErrNegativeAmount = "negative amount"
	// ErrInvalidType is raised when the provided type is different than the expected.
	ErrInvalidType = "invalid type for %s: expected %T, received %T"
)

// EventSetWithdrawAddress defines the event data for the SetWithdrawAddress transaction.
type EventSetWithdrawAddress struct {
	Caller            common.Address
	WithdrawerAddress string
}

// EventWithdrawDelegatorRewards defines the event data for the WithdrawDelegatorRewards transaction.
type EventWithdrawDelegatorRewards struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Hash
	Amount           *big.Int
}

// EventWithdrawValidatorRewards defines the event data for the WithdrawValidatorRewards transaction.
type EventWithdrawValidatorRewards struct {
	ValidatorAddress common.Hash
	Commission       *big.Int
}

// EventApproval defines the event data for the authorization Approve transaction.
type EventApproval struct {
	Owner       common.Address
	Spender     common.Address
	Methods     []string
	AllowedList []string
}

// NewMsgSetWithdrawAddress creates a new MsgSetWithdrawAddress instance.
func NewMsgSetWithdrawAddress(args []interface{}) (*distributiontypes.MsgSetWithdrawAddress, common.Address, error) {
	if len(args) != 2 {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	withdrawerAddress, _ := args[1].(string)

	// If the withdrawer address is a hex address, convert it to a bech32 address.
	if common.IsHexAddress(withdrawerAddress) {
		var err error
		withdrawerAddress, err = sdk.Bech32ifyAddressBytes("evmos", common.HexToAddress(withdrawerAddress).Bytes())
		if err != nil {
			return nil, common.Address{}, err
		}
	}

	msg := &distributiontypes.MsgSetWithdrawAddress{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
		WithdrawAddress:  withdrawerAddress,
	}

	return msg, delegatorAddress, nil
}

// NewMsgWithdrawDelegatorReward creates a new MsgWithdrawDelegatorReward instance.
func NewMsgWithdrawDelegatorReward(args []interface{}) (*distributiontypes.MsgWithdrawDelegatorReward, common.Address, error) {
	if len(args) != 2 {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	validatorAddress, _ := args[1].(string)

	msg := &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
		ValidatorAddress: validatorAddress,
	}

	return msg, delegatorAddress, nil
}

// NewMsgWithdrawValidatorCommission creates a new MsgWithdrawValidatorCommission message.
func NewMsgWithdrawValidatorCommission(args []interface{}) (*distributiontypes.MsgWithdrawValidatorCommission, common.Address, error) {
	if len(args) != 1 {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	msg := &distributiontypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: validatorAddress,
	}

	validatorHexAddr, err := HexAddressFromBech32String(msg.ValidatorAddress)
	if err != nil {
		return nil, common.Address{}, err
	}

	return msg, validatorHexAddr, nil
}

// HexAddressFromBech32String converts a hex address to a bech32 encoded address.
func HexAddressFromBech32String(addr string) (res common.Address, err error) {
	if strings.Contains(addr, sdk.PrefixValidator) {
		valAddr, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			return res, err
		}
		return common.BytesToAddress(valAddr.Bytes()), nil
	}
	return common.BytesToAddress(sdk.MustAccAddressFromBech32(addr)), nil
}

// NewValidatorDistributionInfoRequest creates a new QueryValidatorDistributionInfoRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorDistributionInfoRequest(args []interface{}) (*distributiontypes.QueryValidatorDistributionInfoRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	return &distributiontypes.QueryValidatorDistributionInfoRequest{
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewValidatorOutstandingRewardsRequest creates a new QueryValidatorOutstandingRewardsRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorOutstandingRewardsRequest(args []interface{}) (*distributiontypes.QueryValidatorOutstandingRewardsRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	return &distributiontypes.QueryValidatorOutstandingRewardsRequest{
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewValidatorCommissionRequest creates a new QueryValidatorCommissionRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorCommissionRequest(args []interface{}) (*distributiontypes.QueryValidatorCommissionRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	return &distributiontypes.QueryValidatorCommissionRequest{
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewValidatorSlashesRequest creates a new QueryValidatorSlashesRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorSlashesRequest(method *abi.Method, args []interface{}) (*distributiontypes.QueryValidatorSlashesRequest, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 4, len(args))
	}

	if _, ok := args[1].(uint64); !ok {
		return nil, fmt.Errorf(ErrInvalidType, "startingHeight", uint64(0), args[1])
	}
	if _, ok := args[2].(uint64); !ok {
		return nil, fmt.Errorf(ErrInvalidType, "endingHeight", uint64(0), args[2])
	}

	var input ValidatorSlashesInput
	if err := method.Inputs.Copy(&input, args); err != nil {
		return nil, fmt.Errorf("error while unpacking args to ValidatorSlashesInput struct: %s", err)
	}

	return &distributiontypes.QueryValidatorSlashesRequest{
		ValidatorAddress: input.ValidatorAddress,
		StartingHeight:   input.StartingHeight,
		EndingHeight:     input.EndingHeight,
		Pagination:       &input.PageRequest,
	}, nil
}

// NewDelegationRewardsRequest creates a new QueryDelegationRewardsRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegationRewardsRequest(args []interface{}) (*distributiontypes.QueryDelegationRewardsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	validatorAddress, _ := args[1].(string)

	return &distributiontypes.QueryDelegationRewardsRequest{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewDelegationTotalRewardsRequest creates a new QueryDelegationTotalRewardsRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegationTotalRewardsRequest(args []interface{}) (*distributiontypes.QueryDelegationTotalRewardsRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 1, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	return &distributiontypes.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
	}, nil
}

// NewDelegatorValidatorsRequest creates a new QueryDelegatorValidatorsRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegatorValidatorsRequest(args []interface{}) (*distributiontypes.QueryDelegatorValidatorsRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 1, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	return &distributiontypes.QueryDelegatorValidatorsRequest{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
	}, nil
}

// NewDelegatorWithdrawAddressRequest creates a new QueryDelegatorWithdrawAddressRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegatorWithdrawAddressRequest(args []interface{}) (*distributiontypes.QueryDelegatorWithdrawAddressRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(ErrInvalidNumberOfArgs, 1, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	return &distributiontypes.QueryDelegatorWithdrawAddressRequest{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
	}, nil
}

// ValidatorDistributionInfo is a struct to represent the key information from
// a ValidatorDistributionInfoResponse.
type ValidatorDistributionInfo struct {
	OperatorAddress string        `abi:"operatorAddress"`
	SelfBondRewards []sdk.DecCoin `abi:"selfBondRewards"`
	Commission      []sdk.DecCoin `abi:"commission"`
}

// ValidatorDistributionInfoOutput is a wrapper for ValidatorDistributionInfo to return in the response.
type ValidatorDistributionInfoOutput struct {
	DistributionInfo ValidatorDistributionInfo `abi:"distributionInfo"`
}

// FromResponse converts a response to a ValidatorDistributionInfo.
func (o *ValidatorDistributionInfoOutput) FromResponse(res *distributiontypes.QueryValidatorDistributionInfoResponse) ValidatorDistributionInfoOutput {
	return ValidatorDistributionInfoOutput{
		DistributionInfo: ValidatorDistributionInfo{
			OperatorAddress: res.OperatorAddress,
			SelfBondRewards: NewDecCoinsResponse(res.SelfBondRewards),
			Commission:      NewDecCoinsResponse(res.Commission),
		},
	}
}

// NewDecCoinsResponse converts a response to an array of DecCoin.
func NewDecCoinsResponse(amount sdk.DecCoins) []sdk.DecCoin {
	// Create a new output for each coin and add it to the output array.
	outputs := make([]sdk.DecCoin, len(amount))
	for i, coin := range amount {
		outputs[i] = sdk.DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount,
		}
	}
	return outputs
}

// ValidatorSlashesInput is a struct to represent the key information
// to perform a ValidatorSlashes query.
type ValidatorSlashesInput struct {
	ValidatorAddress string
	StartingHeight   uint64
	EndingHeight     uint64
	PageRequest      query.PageRequest
}

// DelegationDelegatorReward is a struct to represent the key information from
// a query for the rewards of a delegation to a given validator.
type DelegationDelegatorReward struct {
	ValidatorAddress string
	Reward           []sdk.DecCoin
}

// DelegationTotalRewardsOutput is a struct to represent the key information from
// a DelegationTotalRewards response.
type DelegationTotalRewardsOutput struct {
	Rewards []DelegationDelegatorReward
	Total   []sdk.DecCoin
}

// FromResponse populates the DelegationTotalRewardsOutput from a QueryDelegationTotalRewardsResponse.
func (dtr *DelegationTotalRewardsOutput) FromResponse(res *distributiontypes.QueryDelegationTotalRewardsResponse) *DelegationTotalRewardsOutput {
	dtr.Rewards = make([]DelegationDelegatorReward, len(res.Rewards))
	for i, r := range res.Rewards {
		dtr.Rewards[i] = DelegationDelegatorReward{
			ValidatorAddress: r.ValidatorAddress,
			Reward:           NewDecCoinsResponse(r.Reward),
		}
	}
	dtr.Total = NewDecCoinsResponse(res.Total)
	return dtr
}

// Pack packs a given slice of abi arguments into a byte array.
func (dtr *DelegationTotalRewardsOutput) Pack(args abi.Arguments) ([]byte, error) {
	return args.Pack(dtr.Rewards, dtr.Total)
}
