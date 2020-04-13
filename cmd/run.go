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
	// "fmt"

	"fmt"

	impl "github.com/cloudor-io/cloudctl/pkg/Impl"
	"github.com/spf13/cobra"
)

var (
	runArgs impl.RunArgs
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an docker image on the cloud",
	Long: `Usage:
	cloudctl run [OPTIONS] IMAGE [COMMAND] [ARG...]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, token, err := GetLoginToken()
		if err != nil {
			return fmt.Errorf("Error getting user credentails, please log in.")
		}

		runArgs.Args = args
		runEngine, err := impl.NewRunEngine(&runArgs)
		if err != nil {
			return err
		}
		return runEngine.Run(username, token)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().StringVarP(&runArgs.File, "file", "f", "", "job config yaml file")
	runCmd.Flags().StringVarP(&runArgs.Name, "name", "n", "", "job name")
	runCmd.Flags().StringVarP(&runArgs.Region, "region", "", "us-west-2", "region code in the vendor")
	runCmd.Flags().StringVarP(&runArgs.Vendor, "vendor", "", "aws", "cloud vendor name: aws")
	runCmd.Flags().FloatVarP(&runArgs.TimeoutInMin, "timeout", 30.0, "job timeout in minutes")
	runCmd.Flags().StringVarP(&runArgs.InstanceType, "instance-type", "", "aws", "instance-type in the cloud vendor")
	runCmd.Flags().IntVarP(&runArgs.NumInstances, "num-instances", "", 1, "number of instances to launch")
	runCmd.Flags().StringVarP(&runArgs.Input, "input", "i", "", "input, local directory. Use yaml file if use cloud storage for input")
	runCmd.Flags().StringVarP(&runArgs.Output, "output", "o", "", "output, local directory. User yaml file if use cloud storage for output")
	runCmd.Flags().StringVarP(&runArgs.InputMount, "input-mount", "", "", "input directory mounted to the docker image")
	runCmd.Flags().StringVarP(&runArgs.OutputMount, "output-mount", "", "", "output directory mounted the docker image")
}
