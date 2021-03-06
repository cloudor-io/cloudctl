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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
	"github.com/spf13/cobra"
)

var IsLetterOrNumber = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`).MatchString

func signupForm() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Choose your username (length at least 5, alphabetic letters and numbers only): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if len(username) < 5 {
		return "", "", errors.New("Username not long enough.")
	}

	if !IsLetterOrNumber(username) {
		return "", "", errors.New("Username can only contain alphabetic letters and numbers.")
	}

	passwd, err := utils.GetFirstPassword()
	if err != nil {
		fmt.Printf("Error getting password: %v.", err)
		return "", "", err
	}
	passwd2, err := utils.GetRetypePassword(passwd)
	if err != nil {
		fmt.Printf("Error getting password: %v.", err)
		return "", "", err
	}
	if passwd != passwd2 {
		return "", "", errors.New("Typed passwords do not match.")
	}
	return strings.TrimSpace(username), passwd, nil
}

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for the cloudor service",
	Long:  `Sign up for the cloudor service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, passwd, err := signupForm()
		if err != nil {
			return err
		}
		signUpRequest := request.SignupRequest{
			UserName: username,
			Password: passwd,
		}
		signUpBytes, err := json.Marshal(signUpRequest)
		if err != nil {
			return fmt.Errorf("error marshalling signup struct: %v", err)
		}
		resp, err := request.PostCloudor(signUpBytes, nil, nil, "/register")
		if err != nil {
			log.Printf("error posting to server: %v", err)
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			if len(resp.Body()) != 0 {
				return errors.New("remote API error response: " + string(resp.Body()))
			}
			return errors.New("remote API error code " + strconv.Itoa(resp.StatusCode()))
		}
		log.Printf("register successful: %s.\n", string(resp.Body()))
		return nil
	},
}

func init() {
	userCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
