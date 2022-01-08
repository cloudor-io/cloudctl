package impl

import (
	"encoding/json"
	"log"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/spf13/cobra"
)

// ListClusters
func ListClusters(userName, token *string) *[]request.ListClusterResponse {
	resp, err := request.GetCloudor(userName, token, "/cluster")
	cobra.CheckErr(err)

	clusters := []request.ListClusterResponse{}
	err = json.Unmarshal(resp, &clusters)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		cobra.CheckErr(err)
	}
	return &clusters
}
