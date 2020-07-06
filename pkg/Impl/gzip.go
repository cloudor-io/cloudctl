package impl

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"os"
)

func zipFile(src string, dst *os.File) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	// Create a Reader and use ReadAll to get all the bytes from the file.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Write compressed data.
	w := gzip.NewWriter(dst)
	w.Write(content)
	defer w.Close()

	return nil
}
