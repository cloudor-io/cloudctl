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
	"runtime"

	impl "github.com/cloudor-io/cloudctl/pkg/Impl"
	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var (
	listArch string
	listOS   string
)

func needUpdate(username, token *string) (bool, error) {
	latestRelease, err := impl.GetUpdate(username, token, runtime.GOOS, runtime.GOARCH, "latest")
	if err != nil {
		return false, fmt.Errorf("error getting update for OS %s and Arch %s", runtime.GOOS, runtime.GOARCH)
	}
	if len(*latestRelease) == 0 {
		fmt.Printf("could not find the latest release")
		return false, nil
	}

	filename, err := osext.Executable()
	if err != nil {
		return false, err
	}

	myMD5, err := request.GetMD5(filename)
	if err != nil {
		return false, err
	}
	if myMD5 != ((*latestRelease)[0].MD5) {
		return true, nil
	}
	fmt.Printf("current cloudor at %s is the latest release", filename)
	return false, nil
}

// listCmd represents the list command
var updateListCmd = &cobra.Command{
	Use:   "list",
	Short: "list supported architecture and OS",
	Long:  `List supported architecture and OS`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, token := utils.GetLoginToken()
		releases, err := impl.ListUpdates(username, token)
		if err != nil {
			return err
		}
		impl.NewTableView().ViewUpdates(releases)
		need, err := needUpdate(username, token)
		if need {
			fmt.Printf("found new release, please run \"cloudor update\"")
		}
		return err
	},
}

func init() {
	updateCmd.AddCommand(updateListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateListCmd.Flags().StringVarP(&listArch, "arch", "", "", "specify arch")
	updateListCmd.Flags().StringVarP(&listOS, "os", "", "", "specify os")
}
