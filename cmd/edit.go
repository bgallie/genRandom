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
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type templateArgs struct {
	Editor     string
	ConfigFile string
}

var commandArgs []string

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the current OTP Server configuration settings.",
	Long:  "Edit the current OTP Server configuration settings by launching the configured editor in a terminal.",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Flags().Parse(args)
		CheckError(err)
		commandArgs = getEditCommand()
		cmdLine := exec.Command(commandArgs[0], commandArgs[1:]...)
		cmdErr := cmdLine.Run()
		if cmdErr != nil {
			log.Fatalf("Failed to launch editor '%s': %v", strings.Join(commandArgs, " "), cmdErr)
		}
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display the location of the OTP Server configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nOTP Server configuration file: %s\n\n", grViper.ConfigFileUsed())
	},
}

func init() {
	configCmd.AddCommand(editCmd)
	configCmd.AddCommand(infoCmd)
}

// getCommandArgs gets the command and arguments to launch the default editor.
func getEditCommand() []string {
	var err error
	parms := new(templateArgs)
	parms.Editor = grViper.GetString(fmt.Sprintf("editor.%s.editor", osCpuType))
	if parms.Editor == "" {
		log.Fatalf("No editor configured for OS/CPU type '%s'", osCpuType)
	}
	parms.ConfigFile = filepath.Clean(grViper.ConfigFileUsed())
	// command := template.Must(template.New("command").Parse(GetString(fmt.Sprintf("editor.%s.command", osCpuType))))
	editCommand := grViper.GetString(fmt.Sprintf("editor.%s.command", osCpuType))
	command := template.Must(template.New("command").Parse(editCommand))
	var b strings.Builder
	err = command.Execute(&b, parms)
	if err != nil {
		log.Fatal(err)
	}
	commandArgs := strings.Split(b.String(), grViper.GetString("editor.separator"))
	return commandArgs
}
