package impl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"time"

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
	DryRun       bool
	Detach       bool
	TimeoutInMin float64
	NumInstances string
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
	job.Vendors[0].Instances = runArgs.NumInstances
	if runArgs.Input != "" {
		if runArgs.InputMount == "" {
			log.Fatalf("Input mounting point must be specified if input is used.")
		}
		job.Vendors[0].Inputs[0].LocalPath = runArgs.Input
		job.Spec.InputMounts = append(job.Spec.InputMounts, runArgs.InputMount)
	}
	if runArgs.Output != "" {
		if runArgs.OutputMount == "" {
			log.Print("No mount point for output is specified, only getting stdout.")
		}
		job.Vendors[0].Output.LocalPath = runArgs.Output
		job.Spec.OutputMount = runArgs.OutputMount
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

func NewJobByFile(filePath string) (*api.Job, error) {
	job := api.DefaultJob()
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading job yaml file %v.", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, job)
	if err != nil {
		log.Fatalf("Error parsing yaml file %v.", err)
		return nil, err
	}
	return job, nil
}

type RunEngine struct {
	RunArgs *RunArgs
	Job     *api.Job
}

func NewRunEngine(runArgs *RunArgs) (*RunEngine, error) {
	// if file is specified, ignore all other fields except the tag
	runEngine := &RunEngine{
		RunArgs: runArgs,
	}
	if runArgs.File != "" {
		job, err := NewJobByFile(runArgs.File)
		if err != nil {
			return nil, err
		}
		runEngine.Job = job
	} else {
		job := api.DefaultJob()
		err := updateJobByArgs(job, runArgs)
		if err != nil {
			return nil, err
		}
		runEngine.Job = job
	}
	// give it a random name if not specified
	if runArgs.Name == "" {
		runArgs.Name = randomdata.SillyName()
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
		UserName:     *username,
		RunTag:       run.RunArgs.Tag,
		DryRun:       run.RunArgs.DryRun,
		JobName:      run.RunArgs.Name,
		TimeoutInMin: run.RunArgs.TimeoutInMin,
		NumInstances: run.RunArgs.NumInstances,
		YAML:         string(jobBytes),
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
	// somewhere the json is encoded twice, unquote it TODO
	original := string(resp)
	unquoted, err := strconv.Unquote(original)
	if err != nil {
		log.Fatalf("Intenal error while unquoting response: %v", err)
	}
	jobMessage := request.RunJobMessage{}
	err = json.Unmarshal([]byte(unquoted), &jobMessage)
	if err != nil {
		log.Fatalf("Internal error, cann't parse job response: %v", err)
	}
	if runJobRequest.DryRun {
		log.Printf("Dry run successful, estimated cost is %.2f%s", jobMessage.RunInfo.ReservedCredit, jobMessage.RunInfo.RateUnit)
		return nil
	}
	log.Printf("Job submitted successfully.")
	localInput, _ := jobMessage.Job.HasLocals(run.RunArgs.Tag)
	// if no local input, just return. User will poll the job status
	if localInput {
		log.Fatalf("Not implemented")
	}

	if run.RunArgs.Detach {
		log.Printf("Running in detach mode, exiting.")
		return nil
	}
	return run.Wait(&jobMessage)
}

var minPollInterval = 20.0 // seconds

func GetPollInterval(jobMessage *request.RunJobMessage) float64 {
	// Timeout gives some hint about how long the job is
	timeout := jobMessage.RunInfo.TimeoutInMin * 60
	return math.Max(timeout/5.0, minPollInterval)
}

func (run *RunEngine) Wait(jobMessage *request.RunJobMessage) error {
	log.Printf("Sleep for 60 seconds for the instance to boot")
	time.Sleep(60 * time.Second)
	// interval := GetPollInterval(jobMessage)
	// checkPeriod := interval * time.Second
	//
	// ticker := time.NewTicker(checkPeriod)
	// defer func() {
	// 	ticker.Stop()
	// }()
	// count := 0
	// stopped := false
	// for {
	// 	select {
	// 	case <-ticker.C:
	//
	// 	}
	// 	}
	// }
	return nil
}
