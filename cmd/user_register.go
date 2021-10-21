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
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var IsLetterOrNumber = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`).MatchString

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func signupForm() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	errCount := 0
	var email string
	fmt.Print("Please enter your email address: ")
	const MAX_TRIES = 3
	for {
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
		if !valid(email) {
			errCount++
			if errCount == MAX_TRIES {
				return "", "", errors.New("Invalid email address.")
			}
			tries := "tries"
			if MAX_TRIES-errCount == 1 {
				tries = "try"
			}
			fmt.Printf("Email format not correct, please enter your email again (%d %s left);", MAX_TRIES-errCount, tries)
		} else {
			break
		}
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
		return "", "", errors.New("Two passwords do not match.")
	}
	return strings.TrimSpace(email), passwd, nil
}

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for the cloudor service",
	Long:  `Sign up for the cloudor service`,
	Run: func(cmd *cobra.Command, args []string) {
		email, passwd, err := signupForm()
		cobra.CheckErr(err)
		signUpRequest := request.SignupRequest{
			Email:    email,
			Password: passwd,
		}
		signUpBytes, err := json.Marshal(signUpRequest)
		cobra.CheckErr(err)
		resp, err := request.PostCloudor(signUpBytes, nil, nil, "/user/register")
		cobra.CheckErr(err)

		if resp.StatusCode() != http.StatusOK {
			if len(resp.Body()) != 0 {
				cobra.CheckErr(errors.New("remote API error response: " + string(resp.Body())))
			}
			cobra.CheckErr(errors.New("remote API error code " + strconv.Itoa(resp.StatusCode())))
		}
		fmt.Println("User register successfully, please check email to verify the address.")
		viper.Set("default_email", email)
		err = viper.WriteConfig()
		cobra.CheckErr(err)
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
