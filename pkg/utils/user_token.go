package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// GetLoginToken fetches the token
func GetLoginToken() (*string, *string) {
	homeDir, err := os.UserHomeDir()
	CheckErr(err)
	tokenPath := path.Join(homeDir, ".cloudor", ".tokens")
	tokenName := path.Join(tokenPath, ".login")
	f, err := os.Open(tokenName)
	CheckErr(err)
	defer f.Close()
	reader := bufio.NewReader(f)

	userLine, _, err := reader.ReadLine()
	CheckErr(err)
	tokenLine, _, err := reader.ReadLine()
	CheckErr(err)

	userLineTokens := strings.Split(string(userLine), ":")
	tokenLineTokens := strings.Split(string(tokenLine), ":")
	if len(userLineTokens) != 2 || len(tokenLineTokens) != 2 {
		CheckErr(fmt.Errorf("login credentials corrupted. Please login again"))
	}
	if userLineTokens[0] != "user" || tokenLineTokens[0] != "token" {
		CheckErr(fmt.Errorf("login credentials corrupted. Please login again"))
	}
	return &userLineTokens[1], &tokenLineTokens[1]
}
