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
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Pallinder/go-randomdata"
	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/smallfish/simpleyaml"
	"github.com/spf13/cobra"
)

var (
	file string
	name string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  `Create a project with a yaml config file or directly with a docker image`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, token, err := GetLoginToken()
		if err != nil {
			return fmt.Errorf("Error getting user credentails, please log in.")
		}
		// create a backup name if not specified
		if name == "" {
			name = randomdata.SillyName() // TODO lower case
		}
		if file == "" && len(args) == 0 {
			return fmt.Errorf("Either the input file (-f) or the docker image is missing, recommend a name %s", name)
		}

		createInput := request.CreateRequest{
			UserName: *username,
			Name:     name, // add a name anyway
		}
		if file != "" {
			content, err := ioutil.ReadFile(file)
			if err != nil {
				return fmt.Errorf("%v.", err)
			}
			_, err = simpleyaml.NewYaml(content)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			createInput.YAML = string(content)
		} else {
			createInput.Name = name
			createInput.Image = args[0]
		}
		createInputBytes, err := json.Marshal(createInput)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		resp, err := request.PostCloudor(createInputBytes, username, token, "/project/create")

		if err != nil {
			return err
		}
		fmt.Printf("Create project succeeded: %s", string(resp))
		return nil
	},
}

func init() {
	projectCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&file, "file", "f", "", "project config yaml file")
	createCmd.Flags().StringVarP(&name, "name", "n", "", "set project name")
}
