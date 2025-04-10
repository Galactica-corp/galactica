package contracts

import (
	contractutils "github.com/Galactica-corp/galactica/contracts/utils"
	evmtypes "github.com/Galactica-corp/galactica/x/vm/types"
)

func LoadInterchainSenderCallerContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("InterchainSenderCaller.json")
}
