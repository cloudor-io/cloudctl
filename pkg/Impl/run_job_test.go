package impl

import (
	"fmt"
	"os"
	"testing"

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
