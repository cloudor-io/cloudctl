package impl

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/stretchr/testify/assert"
)

func TestNewJobByFile(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	helloWorld := path + "/../../benchmarks/hello_world.yaml"
	job, err := NewJobByFile(helloWorld)
	assert.Equal(t, err, nil, "no error")
	fmt.Printf("job is %+v", *job)
}

func TestUnmarshal(t *testing.T) {
	jobMsgBytes := "{\"user_name\":\"codemk8\",\"id\":\"1berBI2ZhgxVePKg246FtYXEWWD\",\"run_info\":{\"job_name\":\"Racercarnation\",\"instances\":\"1-1\",\"vendor_index\":0,\"created\":1589003262,\"started\":-1,\"finished\":-1,\"last_updated\":1589003262,\"input_stage\":[],\"output_stage\":[]},\"job\":{\"kind\":\"job\",\"version\":\"v1alpha\",\"spec\":{\"image\":\"codemk8/conda_cuda:10.2\",\"envs\":[{\"name\":\"env1\",\"value\":\"1\"}],\"temp\":{}},\"vendors\":[{\"tag\":\"first_choice\",\"name\":\"aws\",\"instance_type\":\"g4dn.xlarge\",\"region\":\"us-west-2\",\"output\":{\"cloud_storage\":{}}}]}}"
	jobMessage := &request.RunJobMessage{}
	err := json.Unmarshal([]byte(jobMsgBytes), jobMessage)
	if err != nil {
		fmt.Printf("err is %v", err)
	} else {
		fmt.Printf("job msg is %+v", *jobMessage)
	}
	assert.Equal(t, err, nil, "no error")
}
