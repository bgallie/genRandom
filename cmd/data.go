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

	"github.com/spf13/cobra"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Generate a stream of (psudo)random bytes.",
	Run: func(cmd *cobra.Command, args []string) {
		generateData(args)
	},
}

func init() {
	rootCmd.AddCommand(dataCmd)
	dataCmd.Flags().StringVarP(&sCnt, "blocks", "", "1", `Write N blocks.`)
	dataCmd.Flags().StringVarP(&sBlock, "bs", "", "512", `Write up to BYTES bytes at a time.
N and BYTES may be followed by the following multiplicative suffixes: c=1, 
w=2, b=512, K=1000, KB=1024, M=1000*1000, MB=1024*1024, G=1000*1000*1000, 
GB=1024*1024*1024, and so on for T, P, E, Z, Y.`)
	myFlagSet = dataCmd.Flags()
	myFlagSet.SetNormalizeFunc(aliasNormalizeFunc)
}

func generateData(args []string) {
	// var err error
	initEngine(args)
	_, err := io.Copy(getOutputFile(), generateRandomStream())
	checkError(err)
	shutdownIkMachine()
}
