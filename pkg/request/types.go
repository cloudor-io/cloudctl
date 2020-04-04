package request

import (
	"github.com/cloudor-io/cloudctl/pkg/api"
)

// LoginResponse defines the body from login request
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

// CreateRequest defines the create request
type CreateRequest struct {
	UserName string `json:"user_name,omitempty"`
	YAML     string `json:"yaml,omitempty"`
	Name     string `json:"name,omitempty"`
	Image    string `json:"image,omitempty"`
}

// RunJobRequest defines the request for running a job
type RunJobRequest struct {
	UserName string `json:"user_name,omitempty"`
	JobName  string `json:"job_name,omitempty"`
	YAML     string `json:"yaml,omitempty"`
}

type JobRunInfo struct {
	// unique id, read-only
	UUID     string `json:"uuid,omitempty"`
	UserName string `json:"user_name,omitempty"`
	// job name, can be auto-generated
	JobName     string           `json:"job_name,omitempty"`
	Vendor      string           `json:"vendor,omitempty"`
	VendorMeta  string           `json:"vendor_meta,omitempty"`
	Created     int64            `json:"created,omitempty"`
	Started     int64            `json:"started,omitempty"`
	Finished    int64            `json:"finished,omitempty"`
	LastUpdated int64            `json:"last_updated,omitempty"`
	InputStage  api.CloudStorage `json:"input_stage,omitempty"`
	OutputStage api.CloudStorage `json:"output_stage,omitempty"`
	Status      string           `json:"status,omitempty"`
}

// Use structured data structure for communication
type RunJobMessage struct {
	UserName string
	RunInfo  JobRunInfo
	Job      api.Job
}
