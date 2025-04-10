package contracts

import (
	contractutils "github.com/Galactica-corp/galactica/contracts/utils"
	evmtypes "github.com/Galactica-corp/galactica/x/vm/types"
)

func LoadDistributionCallerContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("DistributionCaller.json")
}
