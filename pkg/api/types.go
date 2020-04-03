package api

//////////////////////////////////////////////////

type Command struct {
	Entry string
	Args  []string
}

type RunSpec struct {
	Image    string
	Commands []Command
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

	Spec RunSpec
}

//////////////////////////////////////////////////
type StorageSpec struct {
	Type  string // local, s3
	Path  string // local path or cloud path (s3 path e)
	Mount string // mounting path in the docker image
}

// Job
type CloudVendor struct {
	Tag          string
	Name         string
	InstanceType string
	Region       string
	Inputs       []StorageSpec
	Outputs      []StorageSpec
}

type Job struct {
	// must be job
	Kind string
	// unique id, read-only
	UUID string
	// name, can be auto-generated
	Name   string
	RunTag string
	Spec   RunSpec

	Vendors []CloudVendor
}

func DefaultJob() *Job {
	return &Job{
		Kind: "job",
		UUID: "",
		Name: "",
		Spec: RunSpec{
			Image: "",
		},
		RunTag: "first_choice",
		Vendors: []CloudVendor{
			CloudVendor{
				Tag:          "first_choice",
				Name:         "aws",
				InstanceType: "g3s.xlarge",
				Region:       "us-west-2",
				Inputs: []StorageSpec{
					StorageSpec{
						Type: "local",
					},
				},
				Outputs: []StorageSpec{
					StorageSpec{
						Type: "local",
					},
				},
			},
		},
	}
}
