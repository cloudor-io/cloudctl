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
