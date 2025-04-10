package main

import (
	"fmt"
	"os"

	"github.com/Galactica-corp/galactica/cmd/galacticad/cmd"
	evmdconfig "github.com/Galactica-corp/galactica/cmd/galacticad/config"
	examplechain "github.com/Galactica-corp/galactica/galacticad"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	setupSDKConfig()

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "evmd", examplechain.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

func setupSDKConfig() {
	config := sdk.GetConfig()
	evmdconfig.SetBech32Prefixes(config)
	config.Seal()
}
