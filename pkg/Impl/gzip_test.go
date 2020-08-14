package impl

import (
	"fmt"
	"io/ioutil"
	"os"
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

func Test_zipDir(t *testing.T) {
	gzipFile, err := ioutil.TempFile("", "")
	assert.Equal(t, err, nil, "")
	err = zipDir("/tmp/tmp2", gzipFile)
	assert.Equal(t, err, nil, "")
}

func Test_zipDir_relpath(t *testing.T) {
	path, _ := os.Getwd()
	fmt.Printf("cur path is %s", path)
	gzipFile, err := ioutil.TempFile("", "*.zip")
	assert.Equal(t, err, nil, "")
	fmt.Printf("temp file is %s", gzipFile.Name())
	err = zipDir("./tmp2", gzipFile)
	assert.Equal(t, err, nil, "")
}
