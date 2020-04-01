package api

//////////////////////////////////////////////////
// Project
type RunType struct {
	// Supported CPU, NVIDIA-GPU, AMD-GPU, (TPU)
	Type string
	// Minimal core/card
	// optional 1
	Min int32
	// optional -1
	Max int32
	// Minimal Memory in GB
	// optional 4 (GB)
	MinMem int32
}

type Command struct {
	Entry string
	Args  []string
}
type ProjSpec struct {
	Image    string
	Commands []Command
	RunType  RunType
}

type Project struct {
	// must be project
	Kind string

	Name string

	// optional, false
	Public bool
	// optional, free
	Price string
	// optional, empty
	Description string

	Spec ProjSpec
}

//////////////////////////////////////////////////

// Job
type CloudVendor struct {
	Name         string
	InstanceType string
	Region       string
}

type Job struct {
	// must be job
	Kind string
	// unique id, read-only
	UUID    string
	Name    string
	Vendors []CloudVendor
}
