package api

//////////////////////////////////////////////////
type Env struct {
	Name  string `json:"name,omitempty" yaml:"name,string"`
	Value string `json:"value,omitempty" yaml:"value,string"`
}

type TempStorage struct {
	SizeInGB int32  `json:"size_in_gb,omitempty" yaml:"size_in_gb,int"`
	Mount    string `json:"mount,omitempty" yaml:"mount,string"`
}

type RunSpec struct {
	Image string `json:"image,omitempty" yaml:"image,string"`
	// for private docker registry
	ImagePullSecret string      `json:"image_pull_secret,omitempty" yaml:"image_pull_secret,string"`
	Envs            []Env       `json:"envs,omitempty" yaml:"envs,"`
	Command         string      `json:"command,omitempty" yaml:"command"`
	Args            []string    `json:"args,omitempty" yaml:"args"`
	Temp            TempStorage `json:"temp,omitempty" yaml:"temp"`
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
	Entrypoint string `json:"entrypoint,omitempty" yaml:"entry_point,string"`
	Key        string `json:"key,omitempty" yaml:"key,string"`
	Secret     string `json:"secret,omitempty" yaml:"secret,string"`
}

type DataSpec struct {
	Type         string       `json:"type,omitempty" yaml:"type,string"` // local, s3
	Path         string       `json:"path,omitempty" yaml:"path,string"` // local path or cloud path (s3 path e)
	CloudStorage CloudStorage `json:"cloud_storage,omitempty"`
	Mount        string       `json:"mount,omitempty" yaml:"mount,string"` // mounting path in the docker image
}

// Job
type CloudVendor struct {
	Tag          string     `json:"tag,omitempty" yaml:"tag,string"`
	Name         string     `json:"name,omitempty" yaml:"name,string"`
	InstanceType string     `json:"instance_type,omitempty" yaml:"instance_type,string"`
	Region       string     `json:"region,omitempty" yaml:"region,string"`
	Instances    string     `json:"instances,omitempty" yaml:"instances,string"`
	Inputs       []DataSpec `json:"inputs,omitempty" yaml:"inputs"`
	Output       DataSpec   `json:"output,omitempty" yaml:"output"`
}

type Job struct {
	// must be job
	Kind string  `json:"kind,omitempty" yaml:"kind,string"`
	Spec RunSpec `json:"spec,omitempty" yaml:"spec"`

	Vendors []CloudVendor `json:"vendors,omitempty" yaml:"vendors"`
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
				Instances:    "1-1",
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
