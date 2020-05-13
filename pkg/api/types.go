package api

//////////////////////////////////////////////////
type Env struct {
	Name  string `json:"name,omitempty" yaml:"name"`
	Value string `json:"value,omitempty" yaml:"value"`
}

type TempStorage struct {
	SizeInGB int32  `json:"size_in_gb,omitempty" yaml:"size_in_gb"`
	Mount    string `json:"mount,omitempty" yaml:"mount"`
}

type RunSpec struct {
	Type  string `json:"type,omitempty"`
	Image string `json:"image,omitempty" yaml:"image"`
	// for private docker registry
	ImagePullSecret string      `json:"image_pull_secret,omitempty" yaml:"image_pull_secret"`
	Envs            []Env       `json:"envs,omitempty" yaml:"envs"`
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
	Type       string `json:"type,omitempty" yaml:"type"`
	Entrypoint string `json:"entrypoint,omitempty" yaml:"entry_point"`
	Key        string `json:"key,omitempty" yaml:"key"`
	Secret     string `json:"secret,omitempty" yaml:"secret"`
	expiry     int64  `json:"expiry,omitempty" yaml:"expiry"` // set when using cloudor's default stage storage
}

type DataSpec struct {
	LocalPath string `json:"local_path,omitempty" yaml:"local_path"` // local path on the client's machine
	MountPath string `json:"mount_path,omitempty" yaml:"mount_path"` // mounting path in the continaer
	// cloud storage for staging data b/w user and job
	Stage CloudStorage `json:"stage,omitempty" yaml:"stage"`
}

// Job
type CloudVendor struct {
	Tag          string     `json:"tag,omitempty" yaml:"tag"`
	Name         string     `json:"name,omitempty" yaml:"name"`
	InstanceType string     `json:"instance_type,omitempty" yaml:"instance_type"`
	Region       string     `json:"region,omitempty" yaml:"region"`
	Instances    string     `json:"instances,omitempty" yaml:"instances"`
	Inputs       []DataSpec `json:"inputs,omitempty" yaml:"inputs"`
	Output       DataSpec   `json:"output,omitempty" yaml:"output"`
}

type Job struct {
	// must be job
	Kind    string  `json:"kind,omitempty" yaml:"kind"`
	Version string  `json:"version,omitempty" yaml:"version"`
	Spec    RunSpec `json:"spec,omitempty" yaml:"spec"`

	Vendors []CloudVendor `json:"vendors,omitempty" yaml:"vendors"`
}

var DefaultTimeout = 30.0

func DefaultJob() *Job {
	return &Job{
		Kind: "job",
		Spec: RunSpec{
			Type:  "docker",
			Image: "",
		},
		Vendors: []CloudVendor{
			CloudVendor{
				Tag:          "first_choice",
				Name:         "aws",
				InstanceType: "g4dn.xlarge",
				Region:       "us-west-2",
				Instances:    "1-1",
				Inputs: []DataSpec{
					DataSpec{},
				},
				Output: DataSpec{
					LocalPath: "./",
				},
			},
		},
	}
}
