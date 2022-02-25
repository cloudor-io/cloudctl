package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/Pallinder/go-randomdata"
	"github.com/cloudor-io/cloudctl/pkg/api"
	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/cloudor-io/cloudctl/pkg/utils"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
)

type RunArgs struct {
	File string
	Tag  string
	Name string
	// Vendor       string
	// Region       string
	// InstanceType string
	DryRun       bool
	Detach       bool
	TimeoutInMin float64
	NumInstances int
	Input        string
	InputMount   string
	Output       string
	OutputMount  string
	Args         []string
}

// update job info by arguments:
func updateJobByArgs(job *api.Job, runArgs *RunArgs) {
	if len(runArgs.Args) == 0 {
		utils.CheckErr(fmt.Errorf("Please speccify which docker image to run."))
	}
	if _, err := os.Stat(runArgs.Args[0]); err == nil {
		utils.CheckErr(fmt.Errorf("%s looks like a file, forgot -f flag?", runArgs.Args[0]))
	}
	job.Spec.Image = runArgs.Args[0]
	// instances, err := strconv.Atoi(runArgs.NumInstances)
	// utils.CheckErr(err)
	job.Vendors[0].Instances = runArgs.NumInstances
	if runArgs.Input != "" {
		if runArgs.InputMount == "" {
			log.Fatalf("Input mounting point must be specified if input is used.")
		}
		job.Vendors[0].Inputs[0].LocalDir = runArgs.Input
		job.Spec.InputMounts = append(job.Spec.InputMounts, runArgs.InputMount)
	}
	if runArgs.Output != "" {
		if runArgs.OutputMount == "" {
			log.Print("No mount point for output is specified, only getting stdout.")
		}
		job.Vendors[0].Output.LocalDir = runArgs.Output
		job.Spec.OutputMount = runArgs.OutputMount
	}
}

func updateJobByFile(job *api.Job, runArgs *RunArgs) error {
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
		updateJobByArgs(job, runArgs)
		runEngine.Job = job
	}
	// give it a random name if not specified
	if runArgs.Name == "" {
		runArgs.Name = randomdata.SillyName()
	}

	return runEngine, nil
}

func UnmarshalJobMsg(resp *resty.Response) (*api.RunJobMessage, error) {
	// somewhere the json is encoded twice, unquote it TODO
	jobMessage := api.RunJobMessage{}
	err := json.Unmarshal(resp.Body(), &jobMessage)

	if err != nil {
		log.Fatalf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &jobMessage, nil
}

func (run *RunEngine) Run(username, token *string) error {
	jobBytes, err := yaml.Marshal(run.Job)
	if err != nil {
		log.Fatalf("Error marshal job struct to yaml %v", err)
		return err
	}
	//num, err := strconv.Atoi(run.RunArgs.NumInstances)
	//utils.CheckErr(err)
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
	utils.CheckErr(err)
	resp, err := request.PostCloudor(&runJobBytes, username, token, "/job/create")
	utils.CheckErr(err)
	if resp.StatusCode() != http.StatusAccepted && resp.StatusCode() != http.StatusOK {
		if len(resp.Body()) != 0 {
			return errors.New("remote API error response: " + string(resp.Body()))
		}
		return errors.New("remote API error code " + strconv.Itoa(resp.StatusCode()))
	}
	jobMessage, err := UnmarshalJobMsg(resp)
	utils.CheckErr(err)
	if runJobRequest.DryRun {
		fmt.Printf("Dry run successful, estimated cost is %d%s", jobMessage.RunInfo.Cost.ReservedCredit, jobMessage.RunInfo.Cost.RateUnit)
		return nil
	}
	if resp.StatusCode() == http.StatusAccepted {
		err = Upload(jobMessage)
		if err != nil {
			fmt.Printf("error uploading %v", err)
			return err
		}
		runJobBytes, _ := json.Marshal(jobMessage)
		resp, err = request.PostCloudor(&runJobBytes, username, token, "/job/start")
		if err != nil {
			fmt.Printf("Starting job failed %v", err)
			return err
		}
	}
	// vendor := jobMessage.Job.Vendors[*jobMessage.RunInfo.VendorIndex]
	// fmt.Printf("job submitted, running on %s/%s/%s w. timeout %.0f minutes, hourly rate %s%.2f",
	//		vendor.Name, vendor.Region, vendor.InstanceType, jobMessage.RunInfo.TimeoutInMin,
	//		jobMessage.RunInfo.Cost.RateUnit, jobMessage.RunInfo.Cost.HourRate)
	if run.RunArgs.Detach {
		fmt.Println("Running in detach mode, exiting.")
		return nil
	}

	jobMsg, err := CheckingJob(jobMessage, username, token)
	// jobMsg, err := run.Wait(&jobMessage, username, token)
	utils.CheckErr(err)
	return run.Fetch(jobMsg)
}

var minPollInterval = 20.0 // seconds

func GetPollInterval(jobMessage *api.RunJobMessage) int64 {
	// Timeout gives some hint about how long the job is
	timeout := jobMessage.RunInfo.TimeoutInMin * 60
	return int64(math.Max(timeout/5.0, minPollInterval))
}

func (run *RunEngine) Fetch(jobMessage *api.RunJobMessage) error {
	if jobMessage.RunInfo.Stages[len(jobMessage.RunInfo.Stages)-1].Status == "finished" {
		if len(jobMessage.RunInfo.OutputStage) > 0 {
			outputStage := jobMessage.RunInfo.OutputStage[0]
			if outputStage.Cloud.Type == "s3" {
				if outputStage.Pair.Get.URL != "" {
					vendor := jobMessage.Job.Vendors[*jobMessage.RunInfo.VendorIndex]
					if vendor.Output.LocalDir != "" || vendor.Output.Cloud.Type == "" {
						output := "./output.zip"
						if vendor.Output.LocalDir != "" {
							os.MkdirAll(vendor.Output.LocalDir, os.ModePerm)
							output = path.Join(vendor.Output.LocalDir, "output.zip")
						}
						return DownloadFromURL(outputStage.Pair.Get.URL, output)
					}
				}
			}
		}
	} else {
		fmt.Printf("last status not finished, do not fetch output")
	}
	return nil // errors.New("Job returned status not successful.")
}
