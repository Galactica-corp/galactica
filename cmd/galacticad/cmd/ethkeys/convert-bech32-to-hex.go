// Copyright 2024 Galactica Network
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ethkeys

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// ConvertBech32ToHexCmd returns a command to convert a bech32 address to hex
func ConvertBech32ToHexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-bech32-to-hex [bech32address]",
		Short: "Convert a bech32 address to hex and print it",
		Long: `Convert a bech32 address to hex and print it.
Example:
$ galacticad convert-bech32-to-hex gala10jmp6sgh4cc6zt3e8gw05wavvejgr5pwqlc3wn
`,
		Args: cobra.MinimumNArgs(1),
		RunE: runShowCmd,
	}

	return cmd
}

func runShowCmd(cmd *cobra.Command, args []string) (err error) {
	clientCtx := client.GetClientContextFromCmd(cmd)
	clientCtx, err = client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}

	if len(args) != 1 {
		return errors.New("requires address argument")
	}

	bech32Addr := args[0]

	hexAddr, err := ConvertBech32ToHex(bech32Addr)
	if err != nil {
		return err
	}

	fmt.Println(hexAddr)
	return nil
}

func ConvertBech32ToHex(bech32Addr string) (hexAddr string, err error) {
	addr, err := types.AccAddressFromBech32(bech32Addr)
	if err != nil {
		return "", fmt.Errorf("invalid bech32 address: %w", err)
	}

	hexAddr = common.BytesToAddress(addr.Bytes()).String()
	return hexAddr, nil
}
