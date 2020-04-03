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
				Vendors: []CloudVendor{
					CloudVendor{
						Tag:          "",
						Name:         "aws",
						InstanceType: "",
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