package impl

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

// GetTrans returns recent jobs issued by a user
func GetTrans(userName, token *string) (*[]request.TransSchema, error) {
	resp, err := request.GetCloudor(userName, token, "/transaction/user/"+*userName)
	if err != nil {
		log.Printf("getting jobs failed for user %s: %v", *userName, err)
		return nil, err
	}
	// somewhere the json is encoded twice, unquote it TODO
	original := string(resp)
	unquoted, err := strconv.Unquote(original)
	if err != nil {
		log.Printf("Intenal error while unquoting response: %v", err)
		return nil, err
	}
	transactions := []request.TransSchema{}
	err = json.Unmarshal([]byte(unquoted), &transactions)
	if err != nil {
		log.Printf("Internal error, cann't parse transaction response: %v", err)
		return nil, err
	}
	return &transactions, nil
}
