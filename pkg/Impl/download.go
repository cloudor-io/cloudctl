package impl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func DownloadFromURL(url string, filename string) error {
	var w *os.File
	// log.Printf("downloading to file %s", filename)
	if len(filename) > 0 {
		f, err := os.Create(filename)
		if err != nil {
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
	return nil
}

func DownloadSelfFromURL(username, token *string, apiPath string, filename string) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	serverURL := viper.GetString("server") + "/api/v1"
	req, err := http.NewRequest("GET", serverURL+apiPath, nil)
	if username != nil {
		req.Header.Set("From", *username)
	}
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}
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
	return nil
}

func DownloadTmpFromURL(url string) (string, error) {
	w, _ := ioutil.TempFile("", "")

	defer w.Close()
	req, err := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to do GET request, %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get object from url: %d:%s",
			resp.StatusCode, resp.Status)
	}

	if _, err = io.Copy(w, resp.Body); err != nil {
		return "", fmt.Errorf("Failed to write object to file %s, %v", w.Name(), err)
	}
	return w.Name(), nil
}
