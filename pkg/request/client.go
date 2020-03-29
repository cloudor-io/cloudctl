package request

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// ServerURL is the cloudor.io's URL
const ServerURL string = "https://cloudor.io/api/v1"

// PostCloudor issues a POST to ServerURL
func PostCloudor(requestBody []byte, token *string, apiPath string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+*token).
		SetBody(requestBody).
		Post(ServerURL + apiPath)
	if err != nil {
		fmt.Printf("Remote login restful error: %v", err)
		return nil, err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, fmt.Errorf("Unauthorized, please log in.")
	}
	if resp.StatusCode() == http.StatusOK {
		return resp.Body(), nil
	} else {
		body := string(resp.Body())
		if body != "" {
			return nil, fmt.Errorf("Remote API error code %d, response %s.", resp.StatusCode(), body)
		}
		return nil, fmt.Errorf("Remote API error code %d, .", resp.StatusCode())
	}

}

// LoginHandler handles login requets
func LoginCloudor(username, password string) ([]byte, error) {
	client := resty.New()
	response, err := client.R().SetHeader("User-Agent", "CloudCtl").
		SetBasicAuth(username, password).
		Get(ServerURL + "/login")
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
