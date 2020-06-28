package impl

import (
	"fmt"
	"testing"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
)

func TestCheckingJob(t *testing.T) {
	userName, token, err := utils.GetLoginToken()
	if err != nil {
		fmt.Errorf("Error getting user credentails, please log in.")
		return
	}
	jobMsg := request.RunJobMessage{
		UserName: "codemk8",
		ID:       "fake-id",
	}
	CheckingJob(&jobMsg, userName, token)
}
