package testdata

import (
	contractutils "github.com/Galactica-corp/galactica/contracts/utils"
	evmtypes "github.com/Galactica-corp/galactica/x/vm/types"
)

func LoadCounterContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("Counter.json")
}

func LoadCounterFactoryContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("CounterFactory.json")
}
