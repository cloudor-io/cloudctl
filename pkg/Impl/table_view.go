package impl

import (
	"fmt"
	"os"
	"time"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/olekukonko/tablewriter"
)

// View interface
type View interface {
	show(jobs *[]request.RunJobMessage)
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

// View implements the View interface
func (v TableView) View(jobs *[]request.RunJobMessage) {
	data := [][]string{}

	for _, job := range *jobs {
		createdTS := request.GetStatusTs(&job, "created")
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
		vendor := job.Job.Vendors[vendorIndex].Name
		region := job.Job.Vendors[vendorIndex].Region
		instance := job.Job.Vendors[vendorIndex].InstanceType
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
		data = append(data, []string{job.ID, created,
			status, duration, cost, vendor, region, instance})
	}

	v.Table.SetHeader([]string{"ID", "Created", "Status", "Elapsed", "Cost", "Vendor", "Region", "Instance"})
	//v.Table.SetFooter([]string{"", "", "Total", strconv.Itoa(apis.Total)})
	v.Table.SetBorder(true)
	v.Table.AppendBulk(data)
	v.Table.Render()
}
