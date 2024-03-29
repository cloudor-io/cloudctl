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

// creditCmd represents the credit command
var creditCmd = &cobra.Command{
	Use:   "credit",
	Short: "List credit history",
	Long:  `List recent ransaction history`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, token := utils.GetLoginToken()
		transactions, err := impl.GetTrans(username, token)
		if err != nil {
			return err
		}
		impl.NewTableView().ViewTrans(transactions)
		//fmt.Printf("received %d jobs", len(*jobs))

		credit, err := impl.GetCredit(username, token)
		if err != nil {
			fmt.Printf("error getting user credit %v", err)
			return err
		}
		fmt.Printf("Current credit %.2f$, reserved %.2f$", credit.Credit, credit.Reserved)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(creditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// creditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// creditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
