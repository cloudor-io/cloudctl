package request

import (
	"time"

	"github.com/cloudor-io/cloudctl/pkg/api"
)

// LoginResponse defines the body from login request
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type SignupRequest struct {
	UserName string `json:"user_name,omitempty"`
	Password string `json:"password,omitempty"`
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
	DryRun       bool    `json:"dry_run,omitempty" yaml:"dry_run"`
	NumInstances string  `json:"num_instances,omitempty" yaml:"num_instances"`
	TimeoutInMin float64 `json:"timeout_in_min,omitempty" yaml:"timeout_in_min"`
	YAML         string  `json:"yaml,omitempty" yaml:"yaml"`
}

type Cost struct {
	HourRate       float64 `json:"hour_rate,omitempty" yaml:"hour_rate"`
	RateUnit       string  `json:"rate_unit,omitempty" yaml:"rate_unit"`
	ComputeCost    float64 `json:"compute_cost,omitempty" yaml:"compute_cost"`
	EgressGB       float64 `json:"egress_gb,omitempty" yaml:"egress_gb"`
	EgressCost     float64 `json:"egress_cost,omitempty" yaml:"egress_cost"`
	AdjustCost     float64 `json:"adjust_cost,omitempty" yaml:"adjust_cost"`
	ReservedCredit float64 `json:"reserved_credit,omitempty" yaml:"reserved_credit"`
}

type Status struct {
	Code        int32  `json:"return_code,omitempty" yaml:"code"`
	Status      string `json:"status,omitempty" yaml:"status"`
	Description string `json:"description,omitempty" yaml:"description"`
	StdOut      string `json:"std_out,omitempty" yaml:"stdout"`
	UnixTime    int64  `json:"unix_time,omitempty" yaml:"unix_time"`
}

type JobRunInfo struct {
	// unique id, read-only
	// job name, can be auto-generated
	JobName      string             `json:"job_name,omitempty" yaml:"job_name"`
	TimeoutInMin float64            `json:"timeout_in_min,omitempty" yaml:"timeout_in_min"`
	Duration     int64              `json:"duration,omitempty" yaml:"duration"`
	Instances    string             `json:"instances,omitempty" yaml:"instances"`
	Cost         Cost               `json:"cost,omitempty" yaml:"cost"`
	VendorIndex  *int32             `json:"vendor_index,omitempty" yaml:"vendor_index"`
	InputStage   []api.StageStorage `json:"input_stage,omitempty" yaml:"input_stage"`
	OutputStage  []api.StageStorage `json:"output_stage,omitempty" yaml:"output_stage"`
	WorkingDir   string             `json:"working_dir,omitempty" yaml:"working_dir"`
	Stages       []Status           `json:"stages,omitempty" yaml:"stages"`
	// Internal usage
	UpdateNotice api.Notice        `json:"update_notice,omitempty" yaml:"update_notice"`
	Reserved     map[string]string `json:"reserved,omitempty" yaml:"reserved"`
}

// RunJobMessage is the structured data structure for communication
type RunJobMessage struct {
	UserName   string     `json:"user_name,omitempty" yaml:"user_name"`
	Created    int64      `json:"created,omitempty" yaml:"created"`
	ID         string     `json:"id,omitempty" yaml:"id"`
	RunInfo    JobRunInfo `json:"run_info,omitempty" yaml:"run_info"`
	VendorMeta string     `json:"vendor_meta,omitempty" yaml:"vendor_meta"`
	Job        api.Job    `json:"job,omitempty" yaml:"job"`
}

type ListJobRequest struct {
	UserName string `json:"user_name,omitempty"`
	ID       string `json:"id,omitempty"`
	Status   string `json:"status,omitempty"`
}

// JobStatus is used for streaming job status to the client
type JobStatus struct {
	UserName    string `json:"user_name,omitempty"`
	ID          string `json:"id,omitempty"`
	Status      string `json:"status,omitempty"`
	StatusCode  int32  `json:"status_code,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	Description string `json:"description,omitempty"`
}

// AddJobStatus add a stage to job's runtime info
func AddJobStatus(jobMsg *RunJobMessage, status *JobStatus) {
	stage := Status{
		Code:        status.StatusCode,
		Status:      status.Status,
		Description: status.Description,
		UnixTime:    time.Now().Unix(),
	}
	jobMsg.RunInfo.Stages = append(jobMsg.RunInfo.Stages, stage)
}
