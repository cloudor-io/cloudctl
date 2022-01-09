/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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

// listCmd represents the list command
var listClusterCmd = &cobra.Command{
	Use:   "list",
	Short: "List available clusters",
	Long:  `List available clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		username, token := utils.GetLoginToken()
		clusters := impl.ListClusters(username, token)
		fmt.Printf("clusters: %+x", clusters)
	},
}

func init() {
	clusterCmd.AddCommand(listClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
