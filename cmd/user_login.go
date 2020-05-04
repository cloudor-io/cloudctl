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
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh/terminal"
)

// save token to $HOME/.cloudor/.tokens/.login
func saveToken(username *string, token *request.LoginResponse) error {
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Error accessing home directory: %v", err)
		return err
	}
	tokenPath := path.Join(homeDir, ".cloudor", ".tokens")
	err = os.MkdirAll(tokenPath, 0700)
	if err != nil {
		fmt.Printf("Error creating directory: %v.", err)
		return err
	}
	tokenName := path.Join(tokenPath, ".login")
	f, err := os.Create(tokenName)
	if err != nil {
		fmt.Printf("Error creating file %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString("user:" + *username + "\n")
	if err != nil {
		fmt.Printf("Error writing file %v", err)
		return err
	}
	_, err = f.WriteString("token:" + token.Token + "\n")
	if err != nil {
		fmt.Printf("Error writing file %v", err)
		return err
	}
	return nil
}

// GetLoginToken fetches the token
func GetLoginToken() (*string, *string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Error accessing home directory: %v", err)
		return nil, nil, err
	}
	tokenPath := path.Join(homeDir, ".cloudor", ".tokens")
	tokenName := path.Join(tokenPath, ".login")
	f, err := os.Open(tokenName)
	if err != nil {
		fmt.Printf("Error creating file %v", err)
		return nil, nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	userLine, _, err := reader.ReadLine()
	if err != nil {
		fmt.Printf("Error reading file %v", err)
		return nil, nil, err
	}
	tokenLine, _, err := reader.ReadLine()
	if err != nil {
		fmt.Printf("Error reading file %v", err)
		return nil, nil, err
	}
	userLineTokens := strings.Split(string(userLine), ":")
	tokenLineTokens := strings.Split(string(tokenLine), ":")
	if len(userLineTokens) != 2 || len(tokenLineTokens) != 2 {
		return nil, nil, fmt.Errorf("Login credentials corrupted. Please login again.")
	}
	if userLineTokens[0] != "user" || tokenLineTokens[0] != "token" {
		return nil, nil, fmt.Errorf("Login credentials corrupted. Please login again.")
	}
	return &userLineTokens[1], &tokenLineTokens[1], nil
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your username at cloudor: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		fmt.Printf("Error getting password: %v.", err)
		return "", "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to cloudor service",
	Long:  `Login to cloudr service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, password, err := credentials()
		if err != nil {
			return err
		}
		tokenBytes, err := request.LoginCloudor(username, password)
		if err != nil {
			return err
		}
		token := request.LoginResponse{}
		err = json.Unmarshal(tokenBytes, &token)
		if err != nil {
			fmt.Printf("Error parsing tokens: %v, check compatibility.", token)
			return err
		}
		err = saveToken(&username, &token)
		if err == nil {
			fmt.Println("Login succeeded!")
		}
		return err
	},
}

func init() {
	userCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
