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
	UserName string `json:"user_name,omitempty" yaml:"user_name,string"`
	YAML     string `json:"yaml,omitempty" yaml:"yaml,string"`
	Name     string `json:"name,omitempty" yaml:"name,string"`
	Image    string `json:"image,omitempty" yaml:"image,string"`
}

// RunJobRequest defines the request for running a job
type RunJobRequest struct {
	UserName     string  `json:"user_name,omitempty" yaml:"user_name,string"`
	JobName      string  `json:"job_name,omitempty" yaml:"job_name,string"`
	RunTag       string  `json:"run_tag,omitempty" yaml:"run_tag,string"`
	NumInstances int     `json:"num_instances,omitempty" yaml:"num_instances,int"`
	TimeoutInMin float32 `json:"timeout_in_min,omitempty" yaml:"timeout_in_min,"`
	YAML         string  `json:"yaml,omitempty" yaml:"yaml,string"`
}

type JobRunInfo struct {
	// unique id, read-only
	UUID     string `json:"uuid,omitempty" yaml:"uuid,string"`
	UserName string `json:"user_name,omitempty" yaml:"user_name,string"`
	// job name, can be auto-generated
	JobName      string           `json:"job_name,omitempty" yaml:"job_name,string"`
	HourRate     float32          `json:"hour_rate,omitempty" yaml:"hour_rate,string"`
	TimeoutInMin float32          `json:"timeout_in_min,omitempty" yaml:"timeout_in_min"`
	RateUnit     string           `json:"rate_unit,omitempty" yaml:"rate_unit,string"`
	Instances    int32            `json:"instances,omitempty" yaml:"instances,int"`
	Cost         float32          `json:"cost,omitempty" yaml:"cost,float"`
	VendorIndex  int32            `json:"vendor_index" yaml:"vendor_index,int"`
	Created      int64            `json:"created,omitempty" yaml:"created,int"`
	Started      int64            `json:"started,omitempty" yaml:"started,int"`
	Finished     int64            `json:"finished,omitempty" yaml:"finished,int"`
	Duration     int64            `json:"duration,omitempty" yaml:"duration,int"`
	LastUpdated  int64            `json:"last_updated,omitempty" yaml:"last_updated,int"`
	InputStage   api.CloudStorage `json:"input_stage,omitempty" yaml:"input_stage"`
	OutputStage  api.CloudStorage `json:"output_stage,omitempty" yaml:"output_stage"`
	Status       string           `json:"status,omitempty" yaml:"status,string"`
	Reason       string           `json:"reason,omitempty" yaml:"reason,string"`
	StdOut       string           `json:"std_out,omitempty" yaml:"std_out,string"`
}

// Use structured data structure for communication
type RunJobMessage struct {
	UserName   string     `json:"user_name,omitempty" yaml:"user_name,string"`
	RunInfo    JobRunInfo `json:"run_info,omitempty" yaml:"run_info,string"`
	VendorMeta string     `json:"job_meta,omitempty" yaml:"vendor_meta,string"`
	Job        api.Job    `json:"job,omitempty yaml:"job"`
}
