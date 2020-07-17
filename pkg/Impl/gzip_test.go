package impl

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_zipFile(t *testing.T) {
	gzipFile, err := ioutil.TempFile("", "")
	assert.Equal(t, err, nil, "")
	fmt.Printf("temp file is %s", gzipFile.Name())
	err = zipFile("/tmp/image.tar", gzipFile)
	assert.Equal(t, err, nil, "")
}

func Test_zipit(t *testing.T) {
	err := zipit("/tmp/tmp2", "/tmp/tmp2.zip")
	assert.Equal(t, err, nil, "")
}
