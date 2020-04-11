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
	UserName     string `json:"user_name,omitempty"`
	JobName      string `json:"job_name,omitempty"`
	RunTag       string `json:"run_tag,omitempty"`
	NumInstances int    `json:"num_instances,omitempty"`
	YAML         string `json:"yaml,omitempty"`
}

type JobRunInfo struct {
	// unique id, read-only
	UUID     string `json:"uuid,omitempty"`
	UserName string `json:"user_name,omitempty"`
	// job name, can be auto-generated
	JobName     string           `json:"job_name,omitempty"`
	HourRate    float32          `json:"hour_rate,omitempty"`
	RateUnit    string           `json:"rate_unit,omitempty"`
	Instances   int32            `json:"instances,omitempty"`
	Cost        float32          `json:"cost,omitempty"`
	VendorIndex int32            `json:"vendor_index,omitempty"`
	VendorMeta  string           `json:"vendor_meta,omitempty"`
	Created     int64            `json:"created,omitempty"`
	Started     int64            `json:"started,omitempty"`
	Finished    int64            `json:"finished,omitempty"`
	Duration    int64            `json:"duration,omitempty"`
	LastUpdated int64            `json:"last_updated,omitempty"`
	InputStage  api.CloudStorage `json:"input_stage,omitempty"`
	OutputStage api.CloudStorage `json:"output_stage,omitempty"`
	Status      string           `json:"status,omitempty"`
	Reason      string           `json:"reason,omitempty"`
}

// Use structured data structure for communication
type RunJobMessage struct {
	UserName string
	RunInfo  JobRunInfo
	Job      api.Job
}
