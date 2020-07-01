package impl

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFromURL(url string, filename string) error {
	var w *os.File
	// log.Printf("downloading to file %s", filename)
	if len(filename) > 0 {
		f, err := os.Create(filename)
		if err != nil {
			log.Printf("Failed to create the download file %s: %v", filename, err)
			return err
		}
		w = f
	} else {
		w = os.Stdout
	}
	defer w.Close()
	req, err := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to do GET request, %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Failed to get object from url: %d:%s",
			resp.StatusCode, resp.Status)
		log.Print(errMsg)
		return errors.New(errMsg)
	}

	if _, err = io.Copy(w, resp.Body); err != nil {
		log.Printf("Failed to write object to file %s, %v", filename, err)
		return err
	}
	log.Printf("Downloading to file %s succeeded.", filename)
	return nil
}