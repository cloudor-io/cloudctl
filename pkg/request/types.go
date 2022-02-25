package request

import (
	"time"

	"github.com/cloudor-io/cloudctl/pkg/api"
)

// Do not change

// LoginResponse defines the body from login request
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type SignupRequest struct {
	Email    string `json:"email,omitempty"`
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
	NumInstances int     `json:"num_instances,omitempty" yaml:"num_instances"`
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
	ReturnCode  int32  `json:"return_code,omitempty" yaml:"return_code"`
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
	Instances    *int32             `json:"instances,omitempty" yaml:"instances"`
	Cost         Cost               `json:"cost,omitempty" yaml:"cost"`
	VendorIndex  *int32             `json:"vendor_index,omitempty" yaml:"vendor_index"`
	InputStages  []api.StageStorage `json:"input_stages,omitempty" yaml:"input_stages"`
	OutputStage  []api.StageStorage `json:"output_stage,omitempty" yaml:"output_stage"`
	ImageStage   api.StageStorage   `json:"image_stage,omitempty" yaml:"image_stage"`
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
	UserName    string `json:"user_name,omitempty" yaml:"user_name"`
	ID          string `json:"id,omitempty" yaml:"id"`
	Status      string `json:"status,omitempty" yaml:"status"`
	StatusCode  int32  `json:"status_code,omitempty" yaml:"status_code"`
	Vendor      string `json:"vendor,omitempty" yaml:"vendor"`
	Description string `json:"description,omitempty" yaml:"description"`
}

type ListClustersResponse struct {
	Vendor       string   `json:"vendor" yaml:"vendor"`
	Region       string   `json:"region" yaml:"region"`
	InstanceType string   `json:"instance_type" yaml:"instance"`
	Queue        []string `json:"queue,omitempty" yaml:"queue"`
}

// AddJobStatus add a stage to job's runtime info
func AddJobStatus(jobMsg *RunJobMessage, status *JobStatus) {
	stage := Status{
		ReturnCode:  status.StatusCode,
		Status:      status.Status,
		Description: status.Description,
		UnixTime:    time.Now().Unix(),
	}
	jobMsg.RunInfo.Stages = append(jobMsg.RunInfo.Stages, stage)
}

// GetStatusTs returns the unix time in second for a status
func GetStatusTs(jobMsg *RunJobMessage, status string) int64 {
	for _, stage := range jobMsg.RunInfo.Stages {
		if stage.Status == status {
			return stage.UnixTime
		}
	}
	return int64(0)
}

// GetJobDuration gets the elapsed time in second between finished/failed to started
func GetJobDuration(jobMsg *RunJobMessage) int64 {
	startedTS := GetStatusTs(jobMsg, "started")
	if startedTS == int64(0) {
		return int64(-1)
	}
	failedTS := GetStatusTs(jobMsg, "failed")
	if failedTS != int64(0) {
		return failedTS - startedTS
	}
	finishedTS := GetStatusTs(jobMsg, "finished")
	if finishedTS != int64(0) {
		return finishedTS - startedTS
	}
	return int64(-1)
}

// LastStatus gets the last stage's status
func LastStatus(jobMsg *RunJobMessage) string {
	return jobMsg.RunInfo.Stages[len(jobMsg.RunInfo.Stages)-1].Status
}

// CreditSchema defines the credit interface
type CreditSchema struct {
	UserName string  `json:"user_name,omitempty" yaml:"user_name"`
	Credit   float64 `json:"credit,omitempty" yaml:"credit"`
	Reserved float64 `json:"reserved,omitempty" yaml:"reserved"`
	Unit     string  `json:"unit,omitempty" yaml:"unit"`
}

// SupportedOSArch defines a supported os/arch pairs
type SupportedOSArch struct {
	OS      string `json:"os,omitempty" yaml:"os"`
	Arch    string `json:"arch,omitempty" yaml:"arch"`
	Release string `json:"release,omitempty" yaml:"release"`
	MD5     string `json:"md5,omitempty" yaml:"md5"`
}

type JobStat struct {
	Booting         int32
	Running         int32
	AssignInstances int
}

type State struct {
	Idling      int
	Busying     int
	Reserved    int
	Terminating int
	Booting     int
	Max         int
	Tag         string
	JobState    *JobStat
}

type SchedulerStatus struct {
	Waiting int
	Status  map[string]State
}
