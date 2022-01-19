package impl

import (
	"encoding/json"
	"log"

	"github.com/cloudor-io/cloudctl/pkg/api"
	"github.com/cloudor-io/cloudctl/pkg/request"
)

// GetJobs returns recent jobs issued by a user
func GetJobs(userName, token *string) (*[]api.RunJobMessage, error) {
	resp, err := request.GetCloudor(userName, token, "/job/user/"+*userName)
	if err != nil {
		log.Printf("getting jobs failed for user %s: %v", *userName, err)
		return nil, err
	}

	jobMessages := []api.RunJobMessage{}
	err = json.Unmarshal(resp, &jobMessages)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &jobMessages, nil
}
