package api

import "strings"

//////////////////////////////////////////////////
// Find the chosen vendor by tag, if no tag is set, choose the first one
// if no vendor exists, return -1
func (job *Job) FindRunningVendorIndexByTag(tag string) int32 {
	if len(job.Vendors) == 0 {
		return -1
	}
	if tag == "" {
		return 0
	}
	for id, vendor := range job.Vendors {
		if vendor.Tag == tag {
			return int32(id)
		}
	}
	// not found, returns negative
	return -1
}

// HasLocals returns a pair of boolean that indicate if the chosen
// vendor run has any local directory involved.
// They affect the behavior of job run:
// If both false, the job run can immediately return
// If Input is true, the job run needs to wait for the instances to boot to copy the local dir
// If output is true, the job run needs to wait for the run to finish to copy back to local dir
func (job *Job) HasLocals(tag string) (bool, bool) {
	runVendorIndex := job.FindRunningVendorIndexByTag(tag)
	if runVendorIndex < 0 {
		return false, false
	}
	vendor := &job.Vendors[runVendorIndex]
	inputHasLocal, outputHasLocal := false, false
	for _, input := range vendor.Inputs {
		if input.LocalDir != "" {
			inputHasLocal = true
		}
	}
	if vendor.Output.LocalDir != "" {
		outputHasLocal = true
	}
	return inputHasLocal, outputHasLocal
}

// ImageIsLocalFile checks if an image string is a local file or not
func ImageIsLocalFile(image string) bool {
	// Currently just simply check suffix (tar for docker, img for singularity)
	return strings.HasSuffix(image, ".tar") || strings.HasSuffix(image, ".img")
}
