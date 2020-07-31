package request

import (
	"crypto/md5"
	"io"
	"log"
	"os"
)

// GetMD5 returns a md5 of a file
func GetMD5(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return h.Sum(nil), nil
}
