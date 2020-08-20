package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// GetLoginToken fetches the token
func GetLoginToken() (*string, *string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, nil, err
	}
	tokenPath := path.Join(homeDir, ".cloudor", ".tokens")
	tokenName := path.Join(tokenPath, ".login")
	f, err := os.Open(tokenName)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	userLine, _, err := reader.ReadLine()
	if err != nil {
		return nil, nil, err
	}
	tokenLine, _, err := reader.ReadLine()
	if err != nil {
		return nil, nil, err
	}
	userLineTokens := strings.Split(string(userLine), ":")
	tokenLineTokens := strings.Split(string(tokenLine), ":")
	if len(userLineTokens) != 2 || len(tokenLineTokens) != 2 {
		return nil, nil, fmt.Errorf("login credentials corrupted. Please login again")
	}
	if userLineTokens[0] != "user" || tokenLineTokens[0] != "token" {
		return nil, nil, fmt.Errorf("login credentials corrupted. Please login again")
	}
	return &userLineTokens[1], &tokenLineTokens[1], nil
}
