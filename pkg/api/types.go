package api

//////////////////////////////////////////////////
type Env struct {
	Name  string
	Value string
}

type TempStorage struct {
	SizeInGB int32
	Mount    string
}

type RunSpec struct {
	Image   string
	Envs    Env
	Command string
	Args    []string
	Temp    TempStorage
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
type CloudStorage struct {
	Entrypoint string `json:"url,omitempty"`
	Key        string `json:"key,omitempty"`
	Secret     string `json:"secret,omitempty"`
}

type DataSpec struct {
	Type         string       `json:"type,omitempty"` // local, s3
	Path         string       `json:"path,omitempty"` // local path or cloud path (s3 path e)
	CloudStorage CloudStorage `json:"cloud_storage,omitempty"`
	Mount        string       `json:"mount,omitempty"` // mounting path in the docker image
}

// Job
type CloudVendor struct {
	Tag          string     `json:"tag,omitempty"`
	Name         string     `json:"name,omitempty"`
	InstanceType string     `json:"instance_type,omitempty"`
	Region       string     `json:"region,omitempty"`
	Instances    string     `json:"instances,omitempty"`
	Inputs       []DataSpec `json:"inputs,omitempty"`
	Output       DataSpec   `json:"output,omitempty"`
}

type Job struct {
	// must be job
	Kind string  `json:"kind,omitempty"`
	Spec RunSpec `json:"spec,omitempty"`

	Vendors []CloudVendor `json:"vendors,omitempty"`
}

func DefaultJob() *Job {
	return &Job{
		Kind: "job",
		Spec: RunSpec{
			Image: "",
		},
		Vendors: []CloudVendor{
			CloudVendor{
				Tag:          "first_choice",
				Name:         "aws",
				InstanceType: "g3s.xlarge",
				Region:       "us-west-2",
				Instances:    "1-32",
				Inputs: []DataSpec{
					DataSpec{
						Type: "local",
					},
				},
				Output: DataSpec{
					Type: "local",
				},
			},
		},
	}
}
