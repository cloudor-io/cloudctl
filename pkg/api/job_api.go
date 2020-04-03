package api

//////////////////////////////////////////////////
// Find the chosen vendor by tag, if no tag is set, choose the first one
// if no vendor exists, return -1
func (job *Job) FindRunningVendorIndexByTag() int {
	if len(job.Vendors) == 0 {
		return -1
	}
	if job.RunTag == "" {
		return 0
	}
	for id, vendor := range job.Vendors {
		if vendor.Tag == job.RunTag {
			return id
		}
	}
	// not found, returns the first one
	return 0
}

// HasLocals returns a pair of boolean that indicate if the chosen
// vendor run has any local directory involved.
// They affect the behavior of job run:
// If both false, the job run can immediately return
// If Input is true, the job run needs to wait for the instances to boot to copy the local dir
// If output is true, the job run needs to wait for the run to finish to copy back to local dir
func (job *Job) HasLocals() (bool, bool) {
	runVendorIndex := job.FindRunningVendorIndexByTag()
	if runVendorIndex < 0 {
		return false, false
	}
	vendor := &job.Vendors[runVendorIndex]
	inputHasLocal, outputHasLocal := false, false
	for _, input := range vendor.Inputs {
		if input.Type == "local" {
			inputHasLocal = true
		}
	}
	for _, onput := range vendor.Outputs {
		if onput.Type == "local" {
			outputHasLocal = true
		}
	}
	return inputHasLocal, outputHasLocal
}
