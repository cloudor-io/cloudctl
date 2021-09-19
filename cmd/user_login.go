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
	"github.com/spf13/viper"

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
		fmt.Printf("%v.", err)
		return err
	}
	tokenName := path.Join(tokenPath, ".login")
	f, err := os.Create(tokenName)
	if err != nil {
		fmt.Printf("%v", err)
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
	Run: func(cmd *cobra.Command, args []string) {
		username, password, err := credentials()
		cobra.CheckErr(err)
		tokenBytes, err := request.LoginCloudor(username, password, "auth/login")
		cobra.CheckErr(err)
		token := request.LoginResponse{}
		err = json.Unmarshal(tokenBytes, &token)
		cobra.CheckErr(err)
		err = saveToken(&username, &token)
		cobra.CheckErr(err)
		viper.Set("user", username)
		viper.WriteConfig()
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
