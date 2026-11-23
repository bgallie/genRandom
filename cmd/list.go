/*
Copyright © 2025 Billy G. Allie <bill.allie@defiant.mug.org>

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
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List the contents of the OTP Server configuration file.",
	Long:  "List the contents of the OTP Server configuration file in a human-readable format.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nOTP Server Configuration File: %s\n==================================================================================\n",
			grViper.ConfigFileUsed())
		cFile, err := os.Open(grViper.ConfigFileUsed())
		CheckError(err)
		defer cFile.Close()
		_, err = io.Copy(os.Stdout, cFile)
		CheckError(err)
		fmt.Printf("\n==================================================================================\n")
	},
}

func init() {
	configCmd.AddCommand(listCmd)
}
