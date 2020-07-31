package request

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// GetMD5 returns a md5 of a file
func GetMD5(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
