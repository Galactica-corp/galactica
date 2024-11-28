// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package staking

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
)

const (
	// ErrAuthzDoesNotExistOrExpired is raised when the authorization does not exist.
	ErrAuthzDoesNotExistOrExpired = "authorization to %s for address %s does not exist or is expired"
	// ErrEmptyMethods is raised when the given methods array is empty.
	ErrEmptyMethods = "no methods defined; expected at least one message type url"
	// ErrEmptyStringInMethods is raised when the given methods array contains an empty string.
	ErrEmptyStringInMethods = "empty string found in methods array; expected no empty strings to be passed; got: %v"
	// ErrExceededAllowance is raised when the amount exceeds the set allowance.
	ErrExceededAllowance = "amount %s greater than allowed limit %s"
	// ErrInvalidGranter is raised when the granter address is not valid.
	ErrInvalidGranter = "invalid granter address: %v"
	// ErrInvalidGrantee is raised when the grantee address is not valid.
	ErrInvalidGrantee = "invalid grantee address: %v"
	// ErrInvalidMethods is raised when the given methods cannot be unpacked.
	ErrInvalidMethods = "invalid methods defined; expected an array of strings; got: %v"
	// ErrInvalidMethod is raised when the given method cannot be unpacked.
	ErrInvalidMethod = "invalid method defined; expected a string; got: %v"
	// ErrAuthzNotAccepted is raised when the authorization is not accepted.
	ErrAuthzNotAccepted = "authorization to %s for address %s is not accepted"

	// DelegateMethod defines the ABI method name for the staking Delegate
	// transaction.
	DelegateMethod = "delegate"
	// UndelegateMethod defines the ABI method name for the staking Undelegate
	// transaction.
	UndelegateMethod = "undelegate"
	// RedelegateMethod defines the ABI method name for the staking Redelegate
	// transaction.
	RedelegateMethod = "redelegate"
	// CancelUnbondingDelegationMethod defines the ABI method name for the staking
	// CancelUnbondingDelegation transaction.
	CancelUnbondingDelegationMethod = "cancelUnbondingDelegation"

	// DelegateAuthz defines the authorization type for the staking Delegate
	DelegateAuthz = stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE
	// UndelegateAuthz defines the authorization type for the staking Undelegate
	UndelegateAuthz = stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_UNDELEGATE
	// RedelegateAuthz defines the authorization type for the staking Redelegate
	RedelegateAuthz = stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_REDELEGATE
	// CancelUnbondingDelegationAuthz defines the authorization type for the staking
	CancelUnbondingDelegationAuthz = stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION
)

var (
	// DelegateMsg defines the authorization type for MsgDelegate
	DelegateMsg = sdk.MsgTypeURL(&stakingtypes.MsgDelegate{})
	// UndelegateMsg defines the authorization type for MsgUndelegate
	UndelegateMsg = sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{})
	// RedelegateMsg defines the authorization type for MsgRedelegate
	RedelegateMsg = sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{})
	// CancelUnbondingDelegationMsg defines the authorization type for MsgCancelUnbondingDelegation
	CancelUnbondingDelegationMsg = sdk.MsgTypeURL(&stakingtypes.MsgCancelUnbondingDelegation{})
)

func (p Precompile) Delegate(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	bondDemon, err := p.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}

	msg, delegatorHexAddr, err := NewMsgDelegate(args, bondDemon)
	if err != nil {
		return nil, err
	}

	p.Logger(ctx).Debug(
		"tx called",
		"method", method.Name,
		"args", fmt.Sprintf(
			"{ delegator_address: %s, validator_address: %s, amount: %s }",
			delegatorHexAddr,
			msg.ValidatorAddress,
			msg.Amount.Amount,
		),
	)

	var (
		stakeAuthz        *stakingtypes.StakeAuthorization
		expiration        *time.Time
		isCallerOrigin    = contract.CallerAddress == origin
		isCallerDelegator = contract.CallerAddress == delegatorHexAddr
	)

	if isCallerDelegator {
		delegatorHexAddr = origin
	} else if origin != delegatorHexAddr {
		return nil, fmt.Errorf(ErrDifferentOriginFromDelegator, origin.String(), delegatorHexAddr.String())
	}

	if !isCallerOrigin {
		stakeAuthz, expiration, err = CheckAuthzAndAllowanceForGranter(ctx, p.AuthzKeeper, contract.CallerAddress, delegatorHexAddr, &msg.Amount, DelegateMsg)
		if err != nil {
			return nil, err
		}
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(&p.stakingKeeper)
	if _, err = msgSrv.Delegate(sdk.WrapSDKContext(ctx), msg); err != nil {
		return nil, err
	}

	if !isCallerOrigin {
		if err := p.UpdateStakingAuthorization(ctx, contract.CallerAddress, delegatorHexAddr, stakeAuthz, expiration, DelegateMsg, msg); err != nil {
			return nil, err
		}
	}

	if isCallerDelegator {
		stateDB.(*statedb.StateDB).SubBalance(contract.CallerAddress, msg.Amount.Amount.BigInt())
	}

	return method.Outputs.Pack(true)
}

func (p Precompile) Undelegate(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	bonddenom, err := p.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}

	msg, delegatorHexAddr, err := NewMsgUndelegate(args, bonddenom)
	if err != nil {
		return nil, err
	}

	p.Logger(ctx).Debug(
		"tx called",
		"method", method.Name,
		"args", fmt.Sprintf(
			"{ delegator_address: %s, validator_address: %s, amount: %s }",
			delegatorHexAddr,
			msg.ValidatorAddress,
			msg.Amount.Amount,
		),
	)

	var (
		stakeAuthz *stakingtypes.StakeAuthorization
		expiration *time.Time
	)

	if !(contract.CallerAddress == delegatorHexAddr) {
		delegatorHexAddr = origin
	} else if origin != delegatorHexAddr {
		return nil, fmt.Errorf(ErrDifferentOriginFromDelegator, origin.String(), delegatorHexAddr.String())
	}

	if !(contract.CallerAddress == origin) {
		stakeAuthz, expiration, err = CheckAuthzAndAllowanceForGranter(ctx, p.AuthzKeeper, contract.CallerAddress, delegatorHexAddr, &msg.Amount, UndelegateMsg)
		if err != nil {
			return nil, err
		}
	}

	res, err := stakingkeeper.NewMsgServerImpl(&p.stakingKeeper).Undelegate(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, err
	}

	if !(contract.CallerAddress == origin) {
		if err := p.UpdateStakingAuthorization(ctx, contract.CallerAddress, delegatorHexAddr, stakeAuthz, expiration, UndelegateMsg, msg); err != nil {
			return nil, err
		}
	}

	return method.Outputs.Pack(res.CompletionTime.UTC().Unix())
}

func (p Precompile) UpdateStakingAuthorization(
	ctx sdk.Context,
	grantee, granter common.Address,
	stakeAuthz *stakingtypes.StakeAuthorization,
	expiration *time.Time,
	messageType string,
	msg sdk.Msg,
) error {
	updatedResponse, err := stakeAuthz.Accept(ctx, msg)
	if err != nil {
		return err
	}

	if updatedResponse.Delete {
		err = p.AuthzKeeper.DeleteGrant(ctx, grantee.Bytes(), granter.Bytes(), messageType)
	} else {
		err = p.AuthzKeeper.SaveGrant(ctx, grantee.Bytes(), granter.Bytes(), updatedResponse.Updated, expiration)
	}

	if err != nil {
		return err
	}
	return nil
}

func CheckAuthzAndAllowanceForGranter(
	ctx sdk.Context,
	authzKeeper authzkeeper.Keeper,
	grantee, granter common.Address,
	amount *sdk.Coin,
	msgURL string,
) (*stakingtypes.StakeAuthorization, *time.Time, error) {
	msgAuthz, expiration := authzKeeper.GetAuthorization(ctx, grantee.Bytes(), granter.Bytes(), msgURL)
	if msgAuthz == nil {
		return nil, nil, fmt.Errorf(ErrAuthzDoesNotExistOrExpired, msgURL, grantee)
	}

	stakeAuthz, ok := msgAuthz.(*stakingtypes.StakeAuthorization)
	if !ok {
		return nil, nil, authz.ErrUnknownAuthorizationType
	}

	if stakeAuthz.MaxTokens != nil && amount.Amount.GT(stakeAuthz.MaxTokens.Amount) {
		return nil, nil, fmt.Errorf(ErrExceededAllowance, amount.Amount, stakeAuthz.MaxTokens.Amount)
	}

	return stakeAuthz, expiration, nil
}
