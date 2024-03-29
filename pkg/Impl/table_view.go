package impl

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/cloudor-io/cloudctl/pkg/api"
	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/olekukonko/tablewriter"
)

// View interface
type View interface {
	show(jobs *[]api.RunJobMessage)
}

// TableView shows the API list nicely
type TableView struct {
	Table *tablewriter.Table
}

// NewTableView creates the Table
func NewTableView() *TableView {
	return &TableView{
		Table: tablewriter.NewWriter(os.Stdout),
	}
}

// GetStatusTs returns the unix time in second for a status
func GetStatusCode(jobMsg *api.RunJobMessage, status string) (int, error) {
	for _, stage := range jobMsg.RunInfo.Stages {
		if stage.Status == status {
			return int(stage.ReturnCode), nil
		}
	}
	return 0, errors.New("Not found")
}

// View implements the View interface
func (v TableView) View(jobs *[]api.RunJobMessage) {
	/*
		data := [][]string{}

			for _, job := range *jobs {
				createdTS := api.GetStatusTs(&job, "created")
				created := "NA" // time.Unix(job.RunInfo.TimeStamps.Created, 0).Format(time.RFC3339)
				if createdTS > 0 {
					created = time.Unix(createdTS, 0).Format(time.RFC3339)
				}

				vendorIndex := 0
				if job.RunInfo.VendorIndex != nil {
					vendorIndex = int(*job.RunInfo.VendorIndex)
				}
				if len(job.Job.Vendors) <= vendorIndex {
					continue
				}
				duration := "NA"
				durationSec := request.GetJobDuration(&job)
				if durationSec > 3600 {
					duration = fmt.Sprintf("%.1fh", float64(durationSec)/3600.0)
				} else if durationSec > 60 {
					duration = fmt.Sprintf("%.1fm", float64(durationSec)/60.0)
				} else if durationSec > 0 {
					duration = fmt.Sprintf("%ds", durationSec)
				}
				cost := fmt.Sprintf("%.2f", job.RunInfo.Cost.ComputeCost+job.RunInfo.Cost.EgressCost+job.RunInfo.Cost.AdjustCost) + job.RunInfo.Cost.RateUnit
				status := "NA"
				if len(job.RunInfo.Stages) > 0 {
					lastStage := job.RunInfo.Stages[len(job.RunInfo.Stages)-1]
					status = lastStage.Status
				}
				returnCode := "NA"
				code, err := GetStatusCode(&job, "returned")
				if err == nil {
					returnCode = strconv.Itoa(code)
				}
				data = append(data, []string{job.ID, created,
					status, returnCode, duration, cost})
			}

			v.Table.SetHeader([]string{"ID", "Created", "Status", "ReturnCode", "Elapsed", "Cost"})
			//v.Table.SetFooter([]string{"", "", "Total", strconv.Itoa(apis.Total)})
			v.Table.SetBorder(true)
			v.Table.AppendBulk(data)
			v.Table.Render()
	*/
}

// View implements the View interface
func (v TableView) ViewTrans(transactions *[]request.TransSchema) {
	data := [][]string{}

	v.Table.SetHeader([]string{"When", "Type", "Amount", "JobID", "Credit Before", "After"})
	//v.Table.SetFooter([]string{"", "", "Total", strconv.Itoa(apis.Total)})
	for _, trans := range *transactions {

		created := time.Unix(trans.TimeStamp, 0).Format(time.RFC3339)
		jobID := "NA"
		if trans.Info.JobID != "" {
			jobID = trans.Info.JobID
		}
		data = append(data, []string{created, trans.Info.Type,
			fmt.Sprintf("%.2f", trans.Amount) + trans.Unit, jobID,
			fmt.Sprintf("%.2f", trans.Info.CreditBefore) + trans.Unit,
			fmt.Sprintf("%.2f", trans.Info.CreditAfter) + trans.Unit,
		})
	}

	v.Table.SetBorder(true)
	v.Table.AppendBulk(data)
	v.Table.Render()
}

// View implements the View interface
func (v TableView) ViewUpdates(releases *[]request.SupportedOSArch) {
	data := [][]string{}
	myOS := runtime.GOOS
	myArch := runtime.GOARCH
	v.Table.SetHeader([]string{"OS", "Arch", "Release"})
	//v.Table.SetFooter([]string{"", "", "Total", strconv.Itoa(apis.Total)})
	for _, release := range *releases {
		releaseName := release.Release
		if myOS == release.OS && myArch == runtime.GOARCH {
			releaseName += "*"
		}
		data = append(data, []string{release.OS, release.Arch, releaseName})
	}

	v.Table.SetBorder(true)
	v.Table.AppendBulk(data)
	v.Table.Render()
}

// ViewVendors tabulates the vendor
func (v TableView) ViewVendors(instances *[]Instance) {
	data := [][]string{}
	v.Table.SetHeader([]string{"Vendor", "Region", "Type", "CPU", "GPU", "Rate (hr)"})
	//v.Table.SetFooter([]string{"", "", "Total", strconv.Itoa(apis.Total)})
	for _, instance := range *instances {

		data = append(data, []string{instance.Vendor, instance.Region, instance.InstanceType,
			strconv.Itoa(int(instance.Instance.CPU)), instance.Instance.GPU,
			// strconv.Itoa(int(instance.Instance.MemoryGiB)) + "GiB",
			fmt.Sprintf("$%.2f", instance.Instance.Price),
		})
	}

	v.Table.SetBorder(true)
	v.Table.AppendBulk(data)
	v.Table.Render()
}
