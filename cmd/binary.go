/*
Copyright © 2021 Billy G. Allie

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"io"

	"github.com/bgallie/filters/binary"
	"github.com/spf13/cobra"
)

// binaryCmd represents the binary command
var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Generate a stream of (psudo)random ASCII '0' and '1' characters.",
	Run: func(cmd *cobra.Command, args []string) {
		generateBinaryData(args)
	},
}

func init() {
	rootCmd.AddCommand(binaryCmd)
	binaryCmd.Flags().StringVarP(&sCnt, "blocks", "", "1", `Write N blocks.`)
	binaryCmd.Flags().StringVarP(&sBlock, "bs", "", "512", `Write up to BYTES bytes at a time.
N and BYTES may be followed by the following multiplicative suffixes: c=1, 
w=2, b=512, kB=1000, K=1024, MB=1000*1000, M=1024*1024, GB=1000*1000*1000, 
G=1024*1024*1024, and so on for T, P, E, Z, Y.`)
	myFlagSet = binaryCmd.Flags()
	myFlagSet.SetNormalizeFunc(aliasNormalizeFunc)
}

func generateBinaryData(args []string) {
	// var err error
	initEngine(args)
	_, err := io.Copy(getOutputFile(), binary.ToBinary(generateRandomStream()))
	checkError(err)
	shutdownTntMachine()
}
