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
	"io/ioutil"
	"log"
	"os"
	"path"
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
		self, err := osext.Executable()
		if err != nil {
			return err
		}
		originalPath := path.Dir(self)
		tmpFile, err := ioutil.TempFile(originalPath, "cloudor_")
		if err != nil {
			// probably the directory is not writable
			tmpFile, _ = ioutil.TempFile(os.TempDir(), "cloudor_")
		}
		err = impl.DownloadSelfFromURL(username, token, "/releases/"+runtime.GOOS+"/"+runtime.GOARCH+"/latest", tmpFile.Name())
		if err != nil {
			log.Printf("error updating cloudor binary: %v", err)
			return err
		}
		err = os.Rename(tmpFile.Name(), self)
		// remove back up file only after succeeded
		if err != nil {
			log.Printf("error %v", err)
			log.Printf("Latest cloudor written to %s, please manually overwrite %s", tmpFile.Name(), self)
		} else {
			log.Printf("cloudor self updated to %s", self)
		}
		return nil
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
