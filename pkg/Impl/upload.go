package impl

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudor-io/cloudctl/pkg/api"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
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
	fmt.Printf("file %s uploaded", filename)
	return nil
}

// Upload image and inputs, if applicable
func Upload(jobMsg *api.RunJobMessage) error {
	// TODO use go routine to paralize them?
	err := UploadImage(jobMsg)
	if err != nil {
		return err
	}
	err = UploadInputs(jobMsg)
	if err != nil {
		return err
	}
	return nil
}

// UploadImage uploads the local image file to a stage area (S3 presigned URL)
func UploadImage(jobMsg *api.RunJobMessage) error {
	// The Put URL is prepared by the server in create step
	if jobMsg.RunInfo.ImageStage.Pair.Put.URL != "" {
		if !fileExists(jobMsg.Job.Spec.Image) {
			return fmt.Errorf("Image file does not exist %s", jobMsg.Job.Spec.Image)
		}
		gzipFile, err := ioutil.TempFile("", "*.gzip")
		if err != nil {
			return err
		}
		zipFile(jobMsg.Job.Spec.Image, gzipFile)
		defer os.Remove(gzipFile.Name())
		return uploadFile(jobMsg.RunInfo.ImageStage.Pair.Put.URL, gzipFile.Name())
	}
	return nil
}

// UploadInputs uploads
func UploadInputs(jobMsg *api.RunJobMessage) error {
	if len(jobMsg.RunInfo.InputStages) == 0 {
		return nil
	}
	vendor := jobMsg.Job.Vendors[*jobMsg.RunInfo.VendorIndex]
	for inputIndex, stage := range jobMsg.RunInfo.InputStages {
		if stage.Cloud.Type == "s3" {
			//if stage.Pair.Key == "" {
			//				return fmt.Errorf("Expect non-empty s3 pair key in input stage %d", inputIndex)
			//	}
			if vendor.Inputs[inputIndex].LocalDir == "" {
				return fmt.Errorf("Expect local dir when s3 pair keys exist for input %d", inputIndex)
			}
			fmt.Printf("uploading local dir %s", vendor.Inputs[inputIndex].LocalDir)
			err := UploadDirToS3(vendor.Inputs[inputIndex].LocalDir, stage.Pair)
			if err != nil {
				fmt.Printf("error uploading %d input to s3 key: %v", inputIndex, err)
				return fmt.Errorf("%v", err)
			}
		}
	}
	return nil
}

func UploadDirToS3(localDir string, s3Pair api.PresignPair) error {
	if !filepath.IsAbs(localDir) {
		return fmt.Errorf("must be absolute path, %s were given", localDir)
	}

	if !dirExists(localDir) {
		return fmt.Errorf("could not find local dir %s to upload", localDir)
	}
	if s3Pair.Put.URL == "" {
		return fmt.Errorf("missing s3 PUT URL")
	}
	zipFile, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return err
	}
	err = zipDir(localDir, zipFile)
	if err != nil {
		return err
	}
	// fmt.Printf("local dir %s zipped to %s", localDir, zipFile.Name())
	defer os.Remove(zipFile.Name())
	return uploadFile(s3Pair.Put.URL, zipFile.Name())
}
