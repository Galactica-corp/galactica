package testdata

import (
	contractutils "github.com/Galactica-corp/galactica/contracts/utils"
	evmtypes "github.com/Galactica-corp/galactica/x/vm/types"
)

func LoadBankCallerContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("BankCaller.json")
}
