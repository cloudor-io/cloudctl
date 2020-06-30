package impl

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

// GetJobs returns recent jobs issued by a user
func GetJobs(userName, token *string) (*[]request.RunJobMessage, error) {
	resp, err := request.GetCloudor(userName, token, "/job/user/"+*userName)
	if err != nil {
		log.Printf("getting jobs failed for user %s: %v", userName, err)
		return nil, err
	}
	// somewhere the json is encoded twice, unquote it TODO
	original := string(resp)
	unquoted, err := strconv.Unquote(original)
	if err != nil {
		log.Printf("Intenal error while unquoting response: %v", err)
		return nil, err
	}
	jobMessages := []request.RunJobMessage{}
	err = json.Unmarshal([]byte(unquoted), &jobMessages)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &jobMessages, nil
}
