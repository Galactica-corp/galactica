package main

import (
	"fmt"
	"os"

	"github.com/Galactica-corp/galactica/cmd/galacticad/cmd"
	galacticadconfig "github.com/Galactica-corp/galactica/cmd/galacticad/config"
	examplechain "github.com/Galactica-corp/galactica/galacticad"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	setupSDKConfig()

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "galacticad", examplechain.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

func setupSDKConfig() {
	config := sdk.GetConfig()
	galacticadconfig.SetBech32Prefixes(config)
	config.Seal()
}
