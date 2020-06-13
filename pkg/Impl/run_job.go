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
		log.Printf("Dry run successful, estimated cost is %.2f%s", jobMessage.RunInfo.Cost.ReservedCredit, jobMessage.RunInfo.Cost.RateUnit)
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
	jobMsg, err := run.Wait(&jobMessage, username, token)
	if err != nil {
		return err
	}
	return run.Fetch(jobMsg)
}

var minPollInterval = 20.0 // seconds

func GetPollInterval(jobMessage *request.RunJobMessage) int64 {
	// Timeout gives some hint about how long the job is
	timeout := jobMessage.RunInfo.TimeoutInMin * 60
	return int64(math.Max(timeout/5.0, minPollInterval))
}

func (run *RunEngine) Wait(jobMessage *request.RunJobMessage, username, token *string) (*request.RunJobMessage, error) {
	log.Printf("Waiting for the instance to boot")
	time.Sleep(60 * time.Second)
	listJobRequest := request.ListJobRequest{
		UserName: *username,
		ID:       jobMessage.ID,
	}
	listJobBytes, err := json.Marshal(listJobRequest)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	interval := GetPollInterval(jobMessage)
	checkPeriod := time.Second * time.Duration(interval)
	log.Printf("polling job status every %d seconds", interval)
	ticker := time.NewTicker(checkPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			resp, err := request.PostCloudor(listJobBytes, username, token, "/job/list")
			if err != nil {
				log.Fatalf("Submitting job failed %v", err)
				return nil, err
			}
			jobs := []request.RunJobMessage{}
			original := string(resp)
			unquoted, err := strconv.Unquote(original)
			if err != nil {
				log.Fatalf("Intenal error while unquoting response: %v", err)
				return nil, err
			}
			err = json.Unmarshal([]byte(unquoted), &jobs)
			if err != nil {
				log.Fatalf("Internal error, cann't parse job response: %v", err)
				return nil, err
			}
			if jobs[0].Status == "finished" || jobs[0].Status == "failed" {
				log.Printf("Job returned status %s, exit", jobs[0].Status)
				return &jobs[0], nil
			} else {
				log.Printf("Job returned status %s, polling", jobs[0].Status)
			}
		}
	}
}

func (run *RunEngine) Fetch(jobMessage *request.RunJobMessage) error {
	if jobMessage.Status == "finished" {
		if len(jobMessage.RunInfo.OutputStage) > 0 {
			outputStage := jobMessage.RunInfo.OutputStage[0]
			if outputStage.Type == "s3" {
				if outputStage.S3Pair.Get.URL != "" {
					vendor := jobMessage.Job.Vendors[*jobMessage.RunInfo.VendorIndex]
					if vendor.Output.LocalPath != "" || vendor.Output.Stage.Type == "" {
						output := "./output.zip"
						if vendor.Output.LocalPath != "" {
							output = vendor.Output.LocalPath + "output.zip"
						}
						return DownloadFromURL(outputStage.S3Pair.Get.URL, output)
					}
				}
			}
		}
	}
	return nil // errors.New("Job returned status not successful.")
}
