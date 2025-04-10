package testdata

import (
	contractutils "github.com/Galactica-corp/galactica/contracts/utils"
	evmtypes "github.com/Galactica-corp/galactica/x/vm/types"
)

func LoadERC20TestCaller() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("ERC20TestCaller.json")
}
