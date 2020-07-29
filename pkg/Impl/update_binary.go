package impl

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

// ListUpdates returns supported updates
func ListUpdates(userName, token *string) (*[]request.SupportedOSArch, error) {
	resp, err := request.GetCloudor(userName, token, "/update")
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
	updates := []request.SupportedOSArch{}
	err = json.Unmarshal([]byte(unquoted), &updates)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &updates, nil
}

// GetUpdate returns supported upda
func GetUpdate(userName, token *string, os, arch string) (*request.SupportedOSArch, error) {
	resp, err := request.GetCloudor(userName, token, "/update/os/"+os+"/arch"+arch)
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
	updates := request.SupportedOSArch{}
	err = json.Unmarshal([]byte(unquoted), &updates)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &updates, nil
}
