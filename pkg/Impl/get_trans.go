package impl

import (
	"encoding/json"
	"fmt"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

// GetTrans returns recent jobs issued by a user
func GetTrans(userName, token *string) (*[]request.TransSchema, error) {
	resp, err := request.GetCloudor(userName, token, "/transaction/user/"+*userName)
	if err != nil {
		fmt.Printf("getting transactions failed for user %s: %v", *userName, err)
		return nil, err
	}

	transactions := []request.TransSchema{}
	err = json.Unmarshal(resp, &transactions)
	if err != nil {
		fmt.Printf("Internal error, cann't parse transaction response: %v", err)
		return nil, err
	}
	return &transactions, nil
}

func GetCredit(userName, token *string) (*request.CreditSchema, error) {
	resp, err := request.GetCloudor(userName, token, "/credit")
	if err != nil {
		fmt.Printf("getting transactions failed for user %s: %v", *userName, err)
		return nil, err
	}
	credit := request.CreditSchema{}
	err = json.Unmarshal(resp, &credit)
	if err != nil {
		fmt.Printf("Internal error, cann't parse transaction response: %v", err)
		return nil, err
	}
	return &credit, nil
}
