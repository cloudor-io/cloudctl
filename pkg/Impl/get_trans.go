package impl

import (
	"encoding/json"
	"log"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

// GetTrans returns recent jobs issued by a user
func GetTrans(userName, token *string) (*[]request.TransSchema, error) {
	resp, err := request.GetCloudor(userName, token, "/transaction/user/"+*userName)
	if err != nil {
		log.Printf("getting jobs failed for user %s: %v", *userName, err)
		return nil, err
	}

	transactions := []request.TransSchema{}
	err = json.Unmarshal(resp, &transactions)
	if err != nil {
		log.Printf("Internal error, cann't parse transaction response: %v", err)
		return nil, err
	}
	return &transactions, nil
}
