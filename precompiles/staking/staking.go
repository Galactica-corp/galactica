// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package staking

import (
	"bytes"
	"embed"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
)

var (
	stakingContractAddress = common.BytesToAddress([]byte{100})
)

type (
	ErrorOutOfGas    = storetypes.ErrorOutOfGas
	ErrorGasOverflow = storetypes.ErrorGasOverflow
)

const (
	// ErrDecreaseAmountTooBig is raised when the amount by which the allowance should be decreased is greater
	// than the authorization limit.
	ErrDecreaseAmountTooBig = "amount by which the allowance should be decreased is greater than the authorization limit: %s > %s"
	// ErrDifferentOriginFromDelegator is raised when the origin address is not the same as the delegator address.
	ErrDifferentOriginFromDelegator = "origin address %s is not the same as delegator address %s"
	// ErrNoDelegationFound is raised when no delegation is found for the given delegator and validator addresses.
	ErrNoDelegationFound = "delegation with delegator %s not found for validator %s"
)

var _ vm.PrecompiledContract = &Precompile{}

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

// Precompile defines the precompiled contract for staking.
type Precompile struct {
	abi                  abi.ABI
	AuthzKeeper          authzkeeper.Keeper
	stakingKeeper        stakingkeeper.Keeper
	addressContract      common.Address
	kvGasConfig          storetypes.GasConfig
	transientKVGasConfig storetypes.GasConfig
}

// TODO
func (p Precompile) RequiredGas(input []byte) uint64 {
	return 0
}

// NewPrecompile creates a new staking Precompile instance as a
// PrecompiledContract interface.
func NewPrecompile(
	stakingKeeper stakingkeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
) (*Precompile, error) {
	abiBz, err := f.ReadFile("abi.json")
	if err != nil {
		return nil, fmt.Errorf("error loading the staking ABI %s", err)
	}

	newAbi, err := abi.JSON(bytes.NewReader(abiBz))
	if err != nil {
		return nil, err
	}

	addressContract := common.BytesToAddress([]byte{100})

	return &Precompile{
		stakingKeeper:        stakingKeeper,
		AuthzKeeper:          authzKeeper,
		abi:                  newAbi,
		addressContract:      addressContract,
		kvGasConfig:          storetypes.KVGasConfig(),
		transientKVGasConfig: storetypes.TransientGasConfig(),
	}, nil
}

func (Precompile) Address() common.Address {
	return stakingContractAddress
}

// Run executes the precompiled contract staking methods defined in the ABI.
func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) (bz []byte, err error) {
	ctx, stateDB, method, initialGas, args, err := p.RunSetup(evm, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, err
	}

	if err := stateDB.Commit(); err != nil {
		return nil, err
	}

	switch method.Name {
	case DelegateMethod:
		bz, err = p.Delegate(ctx, evm.Origin, contract, stateDB, method, args)
	case UndelegateMethod:
		bz, err = p.Undelegate(ctx, evm.Origin, contract, stateDB, method, args)
	}

	if err != nil {
		return nil, err
	}

	cost := ctx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost) {
		return nil, vm.ErrOutOfGas
	}

	return bz, nil
}

func (Precompile) IsTransaction(method string) bool {
	switch method {
	case DelegateMethod,
		UndelegateMethod,
		RedelegateMethod,
		CancelUnbondingDelegationMethod:
		return true
	default:
		return false
	}
}

// Logger returns a precompile-specific logger.
func (p Precompile) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("evm extension", fmt.Sprintf("x/%s", "staking"))
}

func (p Precompile) RunSetup(
	evm *vm.EVM,
	contract *vm.Contract,
	readOnly bool,
	isTransaction func(name string) bool,
) (ctx sdk.Context, stateDB *statedb.StateDB, method *abi.Method, gasConfig storetypes.Gas, args []interface{}, err error) {
	stateDB, ok := evm.StateDB.(*statedb.StateDB)
	if !ok {
		return sdk.Context{}, nil, nil, uint64(0), nil, fmt.Errorf(ErrNotRunInEvm)
	}
	ctx = stateDB.Context()

	methodID := contract.Input[:4]
	// NOTE: this function iterates over the method map and returns
	// the method with the given ID
	method, err = p.abi.MethodById(methodID)
	if err != nil {
		return sdk.Context{}, nil, nil, uint64(0), nil, err
	}

	// return error if trying to write to state during a read-only call
	if readOnly && isTransaction(method.Name) {
		return sdk.Context{}, nil, nil, uint64(0), nil, vm.ErrWriteProtection
	}

	argsBz := contract.Input[4:]
	args, err = method.Inputs.Unpack(argsBz)
	if err != nil {
		return sdk.Context{}, nil, nil, uint64(0), nil, err
	}

	initialGas := ctx.GasMeter().GasConsumed()

	defer HandleGasError(ctx, contract, initialGas, &err)()

	// set the default SDK gas configuration to track gas usage
	// we are changing the gas meter type, so it panics gracefully when out of gas
	ctx = ctx.WithGasMeter(storetypes.NewGasMeter(contract.Gas)).
		WithKVGasConfig(p.kvGasConfig).
		WithTransientKVGasConfig(p.transientKVGasConfig)
	// we need to consume the gas that was already used by the EVM
	ctx.GasMeter().ConsumeGas(initialGas, "creating a new gas meter")

	return ctx, stateDB, method, initialGas, args, nil
}

func HandleGasError(ctx sdk.Context, contract *vm.Contract, initialGas storetypes.Gas, err *error) func() {
	return func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case ErrorOutOfGas:
				// update contract gas
				usedGas := ctx.GasMeter().GasConsumed() - initialGas
				_ = contract.UseGas(usedGas)

				*err = vm.ErrOutOfGas
				// FIXME: add InfiniteGasMeter with previous Gas limit.
				ctx = ctx.WithKVGasConfig(storetypes.GasConfig{}).
					WithTransientKVGasConfig(storetypes.GasConfig{})
			default:
				panic(r)
			}
		}
	}
}
