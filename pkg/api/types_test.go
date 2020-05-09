package api

import (
	"reflect"
	"testing"
)

func TestDefaultJob(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
		{
			name: "default_value",
			want: &Job{
				Kind: "job",
				Spec: RunSpec{
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
							DataSpec{
								Type: "local",
							},
						},
						Output: DataSpec{
							Type: "local",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultJob(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultJob() = %v, want %v", got, tt.want)
			}
		})
	}
}
