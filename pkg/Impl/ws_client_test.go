package impl

import (
	"testing"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"
)

func TestCheckingJob(t *testing.T) {
	userName, token := utils.GetLoginToken()
	jobMsg := request.RunJobMessage{
		UserName: "codemk8",
		ID:       "fake-id",
	}
	CheckingJob(&jobMsg, userName, token)
}
