package impl

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudor-io/cloudctl/pkg/api"
	"github.com/cloudor-io/cloudctl/pkg/request"
	"gopkg.in/yaml.v2"
)

type RunArgs struct {
	File         string
	Tag          string
	Vendor       string
	Region       string
	InstanceType string
	Input        string
	InputMount   string
	Output       string
	OutputMount  string
	Args         []string
}

// update job info by arguments:
func updateJobByArgs(job *api.Job, runArgs *RunArgs) error {
	if len(runArgs.Args) == 0 {
		log.Fatalf("Please speccify which docker image to run.")
	}
	job.Spec.Image = runArgs.Args[0]
	if runArgs.Vendor != "" {
		job.Vendors[0].Name = runArgs.Vendor
	}
	if runArgs.Region != "" {
		job.Vendors[0].Region = runArgs.Region
	}
	if runArgs.InstanceType != "" {
		job.Vendors[0].InstanceType = runArgs.InstanceType
	}
	if runArgs.Input != "" {
		if runArgs.InputMount == "" {
			log.Fatalf("Input mounting point must be specified if input is used.")
		}
		job.Vendors[0].Inputs[0].Path = runArgs.Input
		job.Vendors[0].Inputs[0].Mount = runArgs.InputMount
	}
	if runArgs.Output != "" {
		if runArgs.OutputMount == "" {
			log.Fatalf("Output mounting point must be specified if output is used.")
		}
		job.Vendors[0].Outputs[0].Path = runArgs.Output
		job.Vendors[0].Outputs[0].Mount = runArgs.OutputMount
	}
	return nil
}

func updateJobByFile(job *api.Job, runArgs *RunArgs) error {
	if runArgs.Vendor != "" {
		log.Printf("vendor argument %s ignored, use yaml file.", runArgs.Vendor)
	}
	log.Panic("To be implemented")
	return nil
}

type RunEngine struct {
	Job *api.Job
}

func NewRunEngine(runArgs *RunArgs) (*RunEngine, error) {
	// if file is specified, ignore all other fields except the tag
	job := api.DefaultJob()
	if runArgs.File != "" {
		err := updateJobByFile(job, runArgs)
		if err != nil {
			return nil, err
		}
	} else {
		err := updateJobByArgs(job, runArgs)
		if err != nil {
			return nil, err
		}
	}
	runEngine := &RunEngine{
		Job: job,
	}
	return runEngine, nil
}

func (run *RunEngine) Run(username, token *string) error {
	jobBytes, err := yaml.Marshal(run.Job)
	if err != nil {
		log.Fatalf("Error marshal job struct to yaml %v", err)
		return err
	}

	runJobRequest := request.RunJobRequest{
		UserName: *username,
		YAML:     string(jobBytes),
	}
	runJobBytes, err := json.Marshal(runJobRequest)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	resp, err := request.PostCloudor(runJobBytes, username, token, "/job/create")
	if err != nil {
		log.Fatalf("Submitting job failed %v", err)
		return err
	}
	log.Printf("Submitting job succeeded: %s", *resp)
	localInput, localOutput := run.Job.HasLocals()
	// if no local input, just return. User will poll the job status
	if !localInput && !localOutput {
		return nil
	}
	log.Fatalf("Not implemented")
	return nil
}
