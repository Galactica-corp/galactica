package distribution

import (
	"bytes"
	"embed"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
)

var _ vm.PrecompiledContract = &Precompile{}

var f embed.FS

type (
	ErrorOutOfGas    = storetypes.ErrorOutOfGas
	ErrorGasOverflow = storetypes.ErrorGasOverflow
)

type Precompile struct {
	distributionKeeper   distributionkeeper.Keeper
	abi                  abi.ABI
	AuthzKeeper          authzkeeper.Keeper
	kvGasConfig          storetypes.GasConfig
	transientKVGasConfig storetypes.GasConfig
}

func NewPrecompile(
	distributionKeeper distributionkeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
) (*Precompile, error) {
	abiBz, err := f.ReadFile("abi.json")
	if err != nil {
		return nil, fmt.Errorf("error loading the distribution ABI %s", err)
	}

	newAbi, err := abi.JSON(bytes.NewReader(abiBz))
	if err != nil {
		return nil, err
	}

	return &Precompile{
		abi:                  newAbi,
		AuthzKeeper:          authzKeeper,
		kvGasConfig:          storetypes.KVGasConfig(),
		transientKVGasConfig: storetypes.TransientGasConfig(),
		distributionKeeper:   distributionKeeper,
	}, nil
}


func (p Precompile) Address() common.Address {
	return common.HexToAddress("0x0000000000000000000000000000000000000801")
}

// TODO
func (p Precompile) RequiredGas(input []byte) uint64 {
	if len(input) < 4 {
		return 0
	}

	// minimum
	return 21000
}

// Run executes the precompiled contract distribution methods defined in the ABI.
func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) (bz []byte, err error) {
	ctx, stateDB, method, initialGas, args, err := p.RunSetup(evm, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, err
	}

	switch method.Name {
	// Distribution transactions
	case SetWithdrawAddressMethod:
		bz, err = p.SetWithdrawAddress(ctx, evm.Origin, contract, stateDB, method, args)
	case WithdrawDelegatorRewardsMethod:
		bz, err = p.WithdrawDelegatorRewards(ctx, evm.Origin, contract, stateDB, method, args)
	case WithdrawValidatorCommissionMethod:
		bz, err = p.WithdrawValidatorCommission(ctx, evm.Origin, contract, stateDB, method, args)
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

// IsTransaction checks if the given methodID corresponds to a transaction or query.
//
// Available distribution transactions are:
//   - SetWithdrawAddress
//   - WithdrawDelegatorRewards
//   - WithdrawValidatorCommission
func (Precompile) IsTransaction(methodID string) bool {
	switch methodID {
	case SetWithdrawAddressMethod,
		WithdrawDelegatorRewardsMethod,
		WithdrawValidatorCommissionMethod:
		return true
	default:
		return false
	}
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

	ctx = ctx.WithGasMeter(storetypes.NewGasMeter(contract.Gas)).WithKVGasConfig(p.kvGasConfig).
		WithTransientKVGasConfig(p.transientKVGasConfig)

	ctx.GasMeter().ConsumeGas(initialGas, "creating a new gas meter")

	return ctx, stateDB, method, initialGas, args, nil
}

func HandleGasError(ctx sdk.Context, contract *vm.Contract, initialGas storetypes.Gas, err *error) func() {
	return func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case ErrorOutOfGas:
				usedGas := ctx.GasMeter().GasConsumed() - initialGas
				_ = contract.UseGas(usedGas)

				*err = vm.ErrOutOfGas
				ctx = ctx.WithKVGasConfig(storetypes.GasConfig{}).
					WithTransientKVGasConfig(storetypes.GasConfig{})
			default:
				panic(r)
			}
		}
	}
}
