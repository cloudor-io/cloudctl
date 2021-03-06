package impl

import (
	"encoding/json"
	"log"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

// ListUpdates returns supported updates
func ListUpdates(userName, token *string) (*[]request.SupportedOSArch, error) {
	resp, err := request.GetCloudor(userName, token, "/update")
	if err != nil {
		log.Printf("getting jobs failed for user %s: %v", *userName, err)
		return nil, err
	}
	updates := []request.SupportedOSArch{}
	err = json.Unmarshal(resp, &updates)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &updates, nil
}

// GetUpdate returns supported upda
func GetUpdate(userName, token *string, os, arch, release string) (*[]request.SupportedOSArch, error) {
	resp, err := request.GetCloudor(userName, token, "/update/os/"+os+"/arch/"+arch+"/release/"+release)
	if err != nil {
		log.Printf("getting jobs failed for user %s: %v", *userName, err)
		return nil, err
	}
	updates := []request.SupportedOSArch{}
	err = json.Unmarshal(resp, &updates)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &updates, nil
}
