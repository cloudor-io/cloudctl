/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"runtime"

	impl "github.com/cloudor-io/cloudctl/pkg/Impl"
	"github.com/cloudor-io/cloudctl/pkg/utils"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update cloudor command line tool",
	Long:  `Update the latest cloudor command line release from cloudor server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, token, err := utils.GetLoginToken()
		if err != nil {
			return fmt.Errorf("error getting update: %v", err)
		}
		need, err := needUpdate(username, token)
		if !need {
			return nil
		}
		filename, err := osext.Executable()
		if err != nil {
			return err
		}
		backupFile := filename + "_bak"
		err = os.Rename(filename, backupFile)
		dstFile := filename
		canReplace := true
		if err != nil {
			log.Printf("error backing up current cloudor binary %v", err)
			dstFile = "cloudor"
			canReplace = false
		}
		err = impl.DownloadSelfFromURL(username, token, "/releases/"+runtime.GOOS+"/"+runtime.GOARCH+"/latest", dstFile)
		if err != nil {
			log.Printf("error updating cloudor binary: %v", err)
		} else {
			// remove back up file only after succeeded
			log.Printf("cloudor self updating completed. You may need to make the file %s to be executable.", dstFile)
			if canReplace {
				os.Remove(backupFile)
			}
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
