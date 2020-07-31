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

	impl "github.com/cloudor-io/cloudctl/pkg/Impl"

	"github.com/cloudor-io/cloudctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	listArch string
	listOS   string
)

// listCmd represents the list command
var updateListCmd = &cobra.Command{
	Use:   "list",
	Short: "list supported architecture and OS",
	Long:  `List supported architecture and OS`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, token, err := utils.GetLoginToken()
		if err != nil {
			return fmt.Errorf("error getting user credentails, please log in (cloudor login)")
		}
		releases, err := impl.ListUpdates(username, token)
		if err != nil {
			return err
		}
		impl.NewTableView().ViewUpdates(releases)
		return nil
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
