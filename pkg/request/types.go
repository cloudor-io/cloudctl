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
	UserName string `json:"user_name,omitempty" yaml:"user_name"`
	YAML     string `json:"yaml,omitempty" yaml:"yaml"`
	Name     string `json:"name,omitempty" yaml:"name"`
	Image    string `json:"image,omitempty" yaml:"image"`
}

// RunJobRequest defines the request for running a job
type RunJobRequest struct {
	UserName     string  `json:"user_name,omitempty" yaml:"user_name"`
	JobName      string  `json:"job_name,omitempty" yaml:"job_name"`
	RunTag       string  `json:"run_tag,omitempty" yaml:"run_tag"`
	NumInstances int     `json:"num_instances,omitempty" yaml:"num_instances"`
	TimeoutInMin float32 `json:"timeout_in_min,omitempty" yaml:"timeout_in_min"`
	YAML         string  `json:"yaml,omitempty" yaml:"yaml"`
}

type JobRunInfo struct {
	// unique id, read-only
	// job name, can be auto-generated
	JobName      string           `json:"job_name,omitempty" yaml:"job_name"`
	HourRate     float32          `json:"hour_rate,omitempty" yaml:"hour_rate"`
	TimeoutInMin float32          `json:"timeout_in_min,omitempty" yaml:"timeout_in_min"`
	RateUnit     string           `json:"rate_unit,omitempty" yaml:"rate_unit"`
	Instances    int32            `json:"instances,omitempty" yaml:"instances"`
	ComputeCost  float32          `json:"compute_cost,omitempty" yaml:"compute_cost"`
	AdjustCost   float32          `json:"adjust_cost,omitempty" yaml:"adjust_cost"`
	VendorIndex  int32            `json:"vendor_index" yaml:"vendor_index"`
	Created      int64            `json:"created,omitempty" yaml:"created"`
	Started      int64            `json:"started,omitempty" yaml:"started"`
	Finished     int64            `json:"finished,omitempty" yaml:"finished"`
	Duration     int64            `json:"duration,omitempty" yaml:"duration"`
	LastUpdated  int64            `json:"last_updated,omitempty" yaml:"last_updated"`
	InputStage   api.CloudStorage `json:"input_stage,omitempty" yaml:"input_stage"`
	OutputStage  api.CloudStorage `json:"output_stage,omitempty" yaml:"output_stage"`
	Status       string           `json:"status,omitempty" yaml:"status"`
	Reason       string           `json:"reason,omitempty" yaml:"reason"`
	StdOut       string           `json:"std_out,omitempty" yaml:"std_out"`
}

// Use structured data structure for communication
type RunJobMessage struct {
	UserName string     `json:"user_name,omitempty" yaml:"user_name"`
	Created  int64      `json:"created,omitempty" yaml:"created"`
	UUID     string     `json:"uuid,omitempty" yaml:"uuid"`
	Status   string     `json:"status,omitempty" yaml:"status"`
	RunInfo  JobRunInfo `json:"run_info,omitempty" yaml:"run_info"`
	Job      api.Job    `json:"job,omitempty" yaml:"job"`
}
