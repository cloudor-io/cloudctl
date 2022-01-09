package impl

import (
	"encoding/json"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
)

// ListClusters
func ListClusters(userName, token *string) *[]request.ListClusterResponse {
	resp, err := request.GetCloudor(userName, token, "/cluster")
	utils.CheckErr(err)

	clusters := []request.ListClusterResponse{}
	err = json.Unmarshal(resp, &clusters)
	utils.CheckErr(err)
	return &clusters
}
