package impl

import (
	"encoding/json"
	"log"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
)

// ListClusters
func ListClusters(userName, token *string) *[]request.ListClusterResponse {
	resp, err := request.GetCloudor(userName, token, "/cluster")
	utils.CheckErr(err)

	clusters := []request.ListClusterResponse{}
	err = json.Unmarshal(resp, &clusters)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		utils.CheckErr(err)
	}
	return &clusters
}
