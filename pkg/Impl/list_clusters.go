package impl

import (
	"encoding/json"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
)

// ListClusters
func ListClusters(userName, token *string) *[]request.ListClustersResponse {
	resp, err := request.PostCloudor(nil, userName, token, "/clusters")
	utils.CheckErr(err)

	clusters := []request.ListClustersResponse{}
	err = json.Unmarshal(resp.Body(), &clusters)
	utils.CheckErr(err)
	return &clusters
}
