// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package staking

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
