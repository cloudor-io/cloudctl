package request

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func TrimQuotes(s string) string {
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
func PostCloudor(requestBody *[]byte, username *string, token *string, apiPath string) (*resty.Response, error) {
	serverURL := viper.GetString("server") + "/api/v1"
	client := resty.New()
	request := client.R().SetHeader("Content-Type", "application/json")
	if requestBody != nil {
		request.SetBody(*requestBody)
	}
	if username != nil {
		request.SetHeader("From", *username)
	}
	if token != nil {
		request.SetHeader("Authorization", "Bearer "+*token)
	}
	resp, err := request.Post(serverURL + apiPath)
	return resp, err
}

// GetCloudor issues a Get to ServerURL
func GetCloudor(username *string, token *string, apiPath string) ([]byte, error) {
	serverURL := viper.GetString("server") + "/api/v1"
	client := resty.New()
	request := client.R()
	if username != nil {
		request.SetHeader("From", *username)
	}
	if token != nil {
		request.SetHeader("Authorization", "Bearer "+*token)
	}
	resp, err := request.Get(serverURL + apiPath)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, errors.New("unauthorized, please log in first")
	}
	if resp.StatusCode() == http.StatusOK {
		return resp.Body(), nil
	}
	if len(resp.Body()) != 0 {
		return nil, errors.New("error code " + strconv.Itoa(resp.StatusCode()) + ":" + string(resp.Body()))
	}
	return nil, errors.New("error code " + strconv.Itoa(resp.StatusCode()))
}

// LoginHandler handles login requets
func LoginCloudor(username, password, urlPath string) ([]byte, error) {
	serverURL := viper.GetString("server") + "/api/v1"
	client := resty.New()
	response, err := client.R().SetHeader("User-Agent", "CloudCtl").
		SetBasicAuth(username, password).
		Get(serverURL + urlPath)
	if err != nil {
		return nil, fmt.Errorf("Login failed: %v", err)
	}
	if response.StatusCode() == http.StatusOK {
		return response.Body(), nil
	} else if response.StatusCode() == http.StatusUnauthorized {
		fmt.Printf("username or password error, please try again.")
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("Login failed with code %d", response.StatusCode())
}
