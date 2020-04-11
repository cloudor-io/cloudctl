package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Pallinder/go-randomdata"
	"github.com/cloudor-io/cloudctl/pkg/api"
	"github.com/cloudor-io/cloudctl/pkg/request"
	"gopkg.in/yaml.v2"
)

type RunArgs struct {
	File         string
	Tag          string
	Name         string
	Vendor       string
	Region       string
	InstanceType string
	NumInstances int
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
	job.Vendors[0].Instances = strconv.Itoa(runArgs.NumInstances)
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
		job.Vendors[0].Output.Path = runArgs.Output
		job.Vendors[0].Output.Mount = runArgs.OutputMount
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
	JobName      string
	RunTag       string
	NumInstances int
	Job          *api.Job
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
	// give it a random name if not specified
	if runArgs.Name == "" {
		runArgs.Name = randomdata.SillyName()
	}
	runEngine := &RunEngine{
		JobName:      runArgs.Name,
		RunTag:       runArgs.Tag,
		NumInstances: runArgs.NumInstances,
		Job:          job,
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
		RunTag:   run.RunTag,
		JobName:  "",
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
	jobMessage := &request.RunJobMessage{}
	err = json.Unmarshal([]byte(*resp), jobMessage)
	if err != nil {
		log.Fatalf("Internal error, cann't parse job response.")
	}
	localInput, localOutput := jobMessage.Job.HasLocals(run.RunTag)
	// if no local input, just return. User will poll the job status
	if !localInput && !localOutput {
		return nil
	}
	log.Fatalf("Not implemented")
	return nil
}
