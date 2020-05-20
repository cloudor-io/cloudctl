package request

import (
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

type JobRunInfo struct {
	// unique id, read-only
	// job name, can be auto-generated
	JobName        string             `json:"job_name,omitempty" yaml:"job_name"`
	HourRate       float64            `json:"hour_rate,omitempty" yaml:"hour_rate"`
	TimeoutInMin   float64            `json:"timeout_in_min,omitempty" yaml:"timeout_in_min"`
	ReservedCredit float64            `json:"reserved_credit,omitempty" yaml:"reserved_credit"`
	RateUnit       string             `json:"rate_unit,omitempty" yaml:"rate_unit"`
	Instances      string             `json:"instances,omitempty" yaml:"instances"`
	ComputeCost    float64            `json:"compute_cost,omitempty" yaml:"compute_cost"`
	EgressGB       float64            `json:"egress_gb,omitempty" yaml:"egress_gb"`
	EgressCost     float64            `json:"egress_cost,omitempty" yaml:"egress_cost"`
	AdjustCost     float64            `json:"adjust_cost,omitempty" yaml:"adjust_cost"`
	VendorIndex    *int32             `json:"vendor_index,omitempty" yaml:"vendor_index"`
	Created        int64              `json:"created,omitempty" yaml:"created"`
	Started        int64              `json:"started,omitempty" yaml:"started"`
	Finished       int64              `json:"finished,omitempty" yaml:"finished"`
	Duration       int64              `json:"duration,omitempty" yaml:"duration"`
	LastUpdated    int64              `json:"last_updated,omitempty" yaml:"last_updated"`
	InputStage     []api.StageStorage `json:"input_stage,omitempty" yaml:"input_stage"`
	OutputStage    []api.StageStorage `json:"output_stage,omitempty" yaml:"output_stage"`
	Reason         string             `json:"reason,omitempty" yaml:"reason"`
	WorkingDir     string             `json:"working_dir,omitempty" yaml:"working_dir"`
	Logs           []string           `json:"logs,omitempty" yaml:"logs"`
	StdOut         string             `json:"std_out,omitempty" yaml:"std_out"`
}

// Use structured data structure for communication
type RunJobMessage struct {
	UserName   string     `json:"user_name,omitempty" yaml:"user_name"`
	Created    int64      `json:"created,omitempty" yaml:"created"`
	ID         string     `json:"id,omitempty" yaml:"id"`
	Status     string     `json:"status,omitempty" yaml:"status"`
	RunInfo    JobRunInfo `json:"run_info,omitempty" yaml:"run_info"`
	VendorMeta string     `json:"vendor_meta,omitempty" yaml:"vendor_meta"`
	Job        api.Job    `json:"job,omitempty" yaml:"job"`
}

type ListJobRequest struct {
	UserName string `json:"user_name,omitempty"`
	ID       string `json:"id,omitempty"`
	Status   string `json:"status,omitempty"`
}
