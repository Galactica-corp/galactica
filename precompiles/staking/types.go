// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package staking

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// EventDelegate defines the event data for the staking Delegate transaction.
type EventDelegate struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Hash
	Amount           *big.Int
	NewShares        *big.Int
}

type EventUnbond struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Hash
	Amount           *big.Int
	CompletionTime   *big.Int
}

type EventRedelegate struct {
	DelegatorAddress    common.Address
	ValidatorSrcAddress common.Hash
	ValidatorDstAddress common.Hash
	Amount              *big.Int
	CompletionTime      *big.Int
}

type EventCancelUnbonding struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Hash
	Amount           *big.Int
	CreationHeight   *big.Int
}

func NewMsgDelegate(args []interface{}, denom string) (*stakingtypes.MsgDelegate, common.Address, error) {
	delegatorAddr, validatorAddress, amount, err := checkDelegationUndelegationArgs(args)
	if err != nil {
		return nil, common.Address{}, err
	}

	msg := &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  denom,
			Amount: math.NewIntFromBigInt(amount),
		},
	}

	return msg, delegatorAddr, nil
}

// NewMsgUndelegate creates a new MsgUndelegate instance and does sanity checks
// on the given arguments before populating the message.
func NewMsgUndelegate(args []interface{}, denom string) (*stakingtypes.MsgUndelegate, common.Address, error) {
	delegatorAddr, validatorAddress, amount, err := checkDelegationUndelegationArgs(args)
	if err != nil {
		return nil, common.Address{}, err
	}

	msg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  denom,
			Amount: math.NewIntFromBigInt(amount),
		},
	}

	return msg, delegatorAddr, nil
}

// NewMsgRedelegate creates a new MsgRedelegate instance and does sanity checks
// on the given arguments before populating the message.
func NewMsgRedelegate(args []interface{}, denom string) (*stakingtypes.MsgBeginRedelegate, common.Address, error) {
	if len(args) != 4 {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidNumberOfArgs, 4, len(args))
	}

	delegatorAddr, ok := args[0].(common.Address)
	if !ok || delegatorAddr == (common.Address{}) {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	validatorSrcAddress, ok := args[1].(string)
	if !ok {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidType, "validatorSrcAddress", "string", args[1])
	}

	validatorDstAddress, ok := args[2].(string)
	if !ok {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidType, "validatorDstAddress", "string", args[2])
	}

	amount, ok := args[3].(*big.Int)
	if !ok {
		return nil, common.Address{}, fmt.Errorf(ErrInvalidAmount, args[3])
	}

	msg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    sdk.AccAddress(delegatorAddr.Bytes()).String(), // bech32 formatted
		ValidatorSrcAddress: validatorSrcAddress,
		ValidatorDstAddress: validatorDstAddress,
		Amount: sdk.Coin{
			Denom:  denom,
			Amount: math.NewIntFromBigInt(amount),
		},
	}

	return msg, delegatorAddr, nil
}

// UnbondingDelegationEntry is a struct that contains the information about an unbonding delegation entry.
type UnbondingDelegationEntry struct {
	CreationHeight int64
	CompletionTime int64
	InitialBalance *big.Int
	Balance        *big.Int
}

// UnbondingDelegationOutput is a struct to represent the key information from
// an unbonding delegation response.
type UnbondingDelegationOutput struct {
	Entries []UnbondingDelegationEntry
}

// FromResponse populates the DelegationOutput from a QueryDelegationResponse.
func (do *UnbondingDelegationOutput) FromResponse(res *stakingtypes.QueryUnbondingDelegationResponse) *UnbondingDelegationOutput {
	do.Entries = make([]UnbondingDelegationEntry, len(res.Unbond.Entries))
	for i, entry := range res.Unbond.Entries {
		do.Entries[i] = UnbondingDelegationEntry{
			CreationHeight: entry.CreationHeight,
			CompletionTime: entry.CompletionTime.UTC().Unix(),
			InitialBalance: entry.InitialBalance.BigInt(),
			Balance:        entry.Balance.BigInt(),
		}
	}
	return do
}

// ValidatorInfo is a struct to represent the key information from
// a validator response.
type ValidatorInfo struct {
	OperatorAddress   string   `abi:"operatorAddress"`
	ConsensusPubkey   string   `abi:"consensusPubkey"`
	Jailed            bool     `abi:"jailed"`
	Status            uint8    `abi:"status"`
	Tokens            *big.Int `abi:"tokens"`
	DelegatorShares   *big.Int `abi:"delegatorShares"` // TODO: Decimal
	Description       string   `abi:"description"`
	UnbondingHeight   int64    `abi:"unbondingHeight"`
	UnbondingTime     int64    `abi:"unbondingTime"`
	Commission        *big.Int `abi:"commission"`
	MinSelfDelegation *big.Int `abi:"minSelfDelegation"`
}

type ValidatorOutput struct {
	Validator ValidatorInfo
}

// DefaultValidatorOutput returns a ValidatorOutput with default values.
func DefaultValidatorOutput() ValidatorOutput {
	return ValidatorOutput{
		ValidatorInfo{
			OperatorAddress:   "",
			ConsensusPubkey:   "",
			Jailed:            false,
			Status:            uint8(0),
			Tokens:            big.NewInt(0),
			DelegatorShares:   big.NewInt(0),
			Description:       "",
			UnbondingHeight:   int64(0),
			UnbondingTime:     int64(0),
			Commission:        big.NewInt(0),
			MinSelfDelegation: big.NewInt(0),
		},
	}
}

// FromResponse populates the ValidatorOutput from a QueryValidatorResponse.
func (vo *ValidatorOutput) FromResponse(res *stakingtypes.QueryValidatorResponse) ValidatorOutput {
	return ValidatorOutput{
		Validator: ValidatorInfo{
			OperatorAddress: res.Validator.OperatorAddress,
			ConsensusPubkey: res.Validator.ConsensusPubkey.String(),
			Jailed:          res.Validator.Jailed,
			Status:          uint8(stakingtypes.BondStatus_value[res.Validator.Status.String()]),
			Tokens:          res.Validator.Tokens.BigInt(),
			DelegatorShares: res.Validator.DelegatorShares.BigInt(), // TODO: Decimal
			// TODO: create description type,
			Description:       res.Validator.Description.Details,
			UnbondingHeight:   res.Validator.UnbondingHeight,
			UnbondingTime:     res.Validator.UnbondingTime.UTC().Unix(),
			Commission:        res.Validator.Commission.CommissionRates.Rate.BigInt(),
			MinSelfDelegation: res.Validator.MinSelfDelegation.BigInt(),
		},
	}
}

// ValidatorsInput is a struct to represent the input information for
// the validators query. Needed to unpack arguments into the PageRequest struct.
type ValidatorsInput struct {
	Status      string
	PageRequest query.PageRequest
}

// checkDelegationUndelegationArgs checks the arguments for the delegation and undelegation functions.
func checkDelegationUndelegationArgs(args []interface{}) (common.Address, string, *big.Int, error) {
	if len(args) != 3 {
		return common.Address{}, "", nil, fmt.Errorf(ErrInvalidNumberOfArgs, 3, len(args))
	}

	delegatorAddr, ok := args[0].(common.Address)
	if !ok || delegatorAddr == (common.Address{}) {
		return common.Address{}, "", nil, fmt.Errorf(ErrInvalidDelegator, args[0])
	}

	validatorAddress, ok := args[1].(string)
	if !ok {
		return common.Address{}, "", nil, fmt.Errorf(ErrInvalidType, "validatorAddress", "string", args[1])
	}

	amount, ok := args[2].(*big.Int)
	if !ok {
		return common.Address{}, "", nil, fmt.Errorf(ErrInvalidAmount, args[2])
	}

	return delegatorAddr, validatorAddress, amount, nil
}
