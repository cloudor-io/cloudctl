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
	impl "github.com/cloudor-io/cloudctl/pkg/Impl"
	"github.com/cloudor-io/cloudctl/pkg/utils"
	"github.com/spf13/cobra"
)

// vendorCmd represents the vendor command
var vendorCmd = &cobra.Command{
	Use:   "vendor",
	Short: "List vendor",
	Long: `Usage:
	cloudor vendor [ARG...]
List supported cloud vendors, it accepts up to two arguments to specify vendor name and region
Examples are:
	cloudor vendor  # show all supported instances
	cloudor vendor aws   # show all supported instances in aws
	cloudor vendor azure eastus  # show all supported instances in azure in eastus region
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, token := utils.GetLoginToken()
		vendors, err := impl.GetVendors(username, token)
		if err != nil {
			return err
		}
		selVendor := ""
		if len(args) >= 1 {
			selVendor = args[0]
		}
		selRegion := ""
		if len(args) >= 2 {
			selRegion = args[1]
		}
		vendorArray := impl.FilterVendors(vendors, selVendor, selRegion)
		impl.NewTableView().ViewVendors(&vendorArray)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(vendorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vendorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vendorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
