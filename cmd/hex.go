/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/bgallie/filters/hex"
	"github.com/bgallie/filters/lines"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	myFlagSet *pflag.FlagSet
)

// hexCmd represents the hex command
var hexCmd = &cobra.Command{
	Use:   "hex",
	Short: "Generate a stream of (psudo)random hexadecimal encoded bytes.",
	Run: func(cmd *cobra.Command, args []string) {
		generateHexData(args)
	},
}

func init() {
	rootCmd.AddCommand(hexCmd)
	hexCmd.Flags().StringVarP(&sCnt, "blocks", "", "1", `Write N blocks.`)
	hexCmd.Flags().StringVarP(&sBlock, "bs", "", "512", `Write up to BYTES bytes at a time.
N and BYTES may be followed by the following multiplicative suffixes: c=1,
w=2, b=512, K=1000, KB=1024, M=1000*1000, MB=1024*1024, G=1000*1000*1000,
GB=1024*1024*1024, and so on for T, P, E, Z, Y.`)
	myFlagSet = hexCmd.Flags()
	myFlagSet.SetNormalizeFunc(aliasNormalizeFunc)
}

func generateHexData(args []string) {
	// var err error
	initEngine(args)
	lines.LineSize = 64
	out, err := getOutputFile()
	cobra.CheckErr(err)
	_, err = io.Copy(out, lines.SplitToLines(hex.ToHex(generateRandomStream())))
	checkError(err)
}
