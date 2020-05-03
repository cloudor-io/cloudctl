package request

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// ServerURL is the cloudor.io's URL
// const ServerURL string = "https://cloudor.io/api/v1"
// const ServerURL string = "https://cloudor.dev/api/v1"

// PostCloudor issues a POST to ServerURL
func PostCloudor(requestBody []byte, username *string, token *string, apiPath string) (*string, error) {
	serverURL := viper.GetString("server")
	client := resty.New()
	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+*token).
		SetHeader("From", *username).
		SetBody(requestBody).
		Post(serverURL + apiPath)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, errors.New("unauthorized, please log in.")
	}
	body := strings.TrimSuffix(trimQuotes(string(resp.Body())), "\\n")
	if resp.StatusCode() == http.StatusOK {
		return &body, nil
	} else {
		if body != "" {

			return nil, errors.New("remote API error response: " + body)
		}
		return nil, errors.New("remote API error code " + strconv.Itoa(resp.StatusCode()))
	}

}

// LoginHandler handles login requets
func LoginCloudor(username, password string) ([]byte, error) {
	serverURL := viper.GetString("server")
	client := resty.New()
	response, err := client.R().SetHeader("User-Agent", "CloudCtl").
		SetBasicAuth(username, password).
		Get(serverURL + "/login")
	if err != nil {
		return nil, fmt.Errorf("Login failed: %v", err)
	}
	if response.StatusCode() == http.StatusOK {
		return response.Body(), nil
	} else if response.StatusCode() == http.StatusUnauthorized {
		fmt.Printf("username or password error, please try again.")
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("Login failed with code %d", response.StatusCode())
}
