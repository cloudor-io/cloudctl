package impl

import (
	"encoding/json"
	"fmt"

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

// ListClusters
func GetClusterStatus(userName, token *string) request.SchedulerStatus {
	resp, err := request.PostCloudor(nil, userName, token, "/cluster/status")
	utils.CheckErr(err)

	status := request.SchedulerStatus{}
	err = json.Unmarshal(resp.Body(), &status)
	fmt.Printf("%s\n", string(resp.Body()))
	utils.CheckErr(err)
	return status
}
