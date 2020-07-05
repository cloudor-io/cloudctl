package impl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func uploadFile(presignURL, filename string) error {
	var r io.ReadCloser

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open upload file %s, %v", filename, err)
	}

	// Get the size of the file so that the constraint of Content-Length
	// can be included with the presigned URL. This can be used by the
	// server or client to ensure the content uploaded is of a certain size.
	//
	// These constraints can further be expanded to include things like
	// Content-Type. Additionally constraints such as X-Amz-Content-Sha256
	// header set restricting the content of the file to only the content
	// the client initially made the request with. This prevents the object
	// from being overwritten or used to upload other unintended content.
	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file, %s, %v", filename, err)
	}

	size := stat.Size()
	r = f

	defer r.Close()

	// Get the Presigned URL from the remote service. Pass in the file's
	// size if it is known so that the presigned URL returned will be required
	// to be used with the size of content requested.
	log.Printf("uploading to %s", presignURL)
	req, err := http.NewRequest("PUT", presignURL, nil)

	if err != nil {
		return fmt.Errorf("failed to get put presigned request, %v", err)
	}
	req.ContentLength = size
	req.Body = r

	// Upload the file contents to S3.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do PUT request, %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to put S3 object, %d:%s",
			resp.StatusCode, resp.Status)
	}

	return nil
}

func Upload(jobMsg *request.RunJobMessage) error {
	return UploadImage(jobMsg)
}

func UploadImage(jobMsg *request.RunJobMessage) error {
	if jobMsg.RunInfo.ImageStage.S3Pair.Put.URL != "" {
		if !fileExists(jobMsg.Job.Spec.Image) {
			return fmt.Errorf("Image file does not exist %s", jobMsg.Job.Spec.Image)
		}
		return uploadFile(jobMsg.RunInfo.ImageStage.S3Pair.Put.URL, jobMsg.Job.Spec.Image)
	}
	return nil
}
