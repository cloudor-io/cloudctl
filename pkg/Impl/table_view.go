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
		created := time.Unix(job.RunInfo.TimeStamps.Created, 0).Format(time.RFC3339)

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
		if job.RunInfo.TimeStamps.Duration > 3600 {
			duration = fmt.Sprintf("%.1fh", job.RunInfo.TimeStamps.Duration/3600.0)
		} else if job.RunInfo.TimeStamps.Duration > 60 {
			duration = fmt.Sprintf("%.1fm", job.RunInfo.TimeStamps.Duration/60.0)
		} else if job.RunInfo.TimeStamps.Duration > 0 {
			duration = fmt.Sprintf("%fs", job.RunInfo.TimeStamps.Duration)
		}
		cost := fmt.Sprintf("%.2f", job.RunInfo.Cost.ComputeCost+job.RunInfo.Cost.EgressCost+job.RunInfo.Cost.AdjustCost) + job.RunInfo.Cost.RateUnit

		data = append(data, []string{job.ID, created,
			job.Status, duration, cost, vendor, region, instance})
	}

	v.Table.SetHeader([]string{"ID", "Created", "Status", "Elapsed", "Cost", "Vendor", "Region", "Instance"})
	//v.Table.SetFooter([]string{"", "", "Total", strconv.Itoa(apis.Total)})
	v.Table.SetBorder(true)
	v.Table.AppendBulk(data)
	v.Table.Render()
}
