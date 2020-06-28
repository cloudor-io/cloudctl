package cmd

import (
	"fmt"

	"github.com/cloudor-io/cloudctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	user string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list projects",
	Long:  `list projects from current user (default) or public projects from other users`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _, err := utils.GetLoginToken()
		if err != nil {
			return fmt.Errorf("Error getting user credentails, please log in.")
		}
		// default to myself
		if user == "" {
			user = *username
		}

		fmt.Println("list called")
		return nil
	},
}

func init() {
	projectCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
