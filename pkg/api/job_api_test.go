package api

import (
	"testing"
)

func TestJob_FindRunningVendorIndexByTag(t *testing.T) {

	type fields struct {
		Kind    string
		UUID    string
		Name    string
		Spec    RunSpec
		Vendors []CloudVendor
	}
	exampleField := fields{
		Vendors: []CloudVendor{
			CloudVendor{
				Tag: "tag1",
			},
			CloudVendor{
				Tag: "tag2",
			},
		},
	}
	tests := []struct {
		name   string
		fields fields
		runTag string
		want   int
	}{
		// TODO: Add test cases.
		{
			name:   "no_vendor_should_return_negative",
			fields: fields{},
			runTag: "",
			want:   -1,
		},
		{
			name:   "no_tag_should_return_index_0",
			fields: exampleField,
			runTag: "",
			want:   0,
		},
		{
			name:   "tag2_should_return_index_1",
			fields: exampleField,
			runTag: "tag2",
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				Kind:    tt.fields.Kind,
				Spec:    tt.fields.Spec,
				Vendors: tt.fields.Vendors,
			}
			if got := job.FindRunningVendorIndexByTag(tt.runTag); got != tt.want {
				t.Errorf("Job.FindRunningVendorIndexByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_HasLocals(t *testing.T) {
	type fields struct {
		Kind    string
		UUID    string
		Name    string
		Spec    RunSpec
		Vendors []CloudVendor
		RunTag  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
		want1  bool
	}{
		// TODO: Add test cases.
		{
			name:   "no_vendor_should_return_falses",
			fields: fields{},
			want:   false,
			want1:  false,
		},
		{
			name: "no_vendor_should_return_falses",
			fields: fields{
				Vendors: []CloudVendor{
					CloudVendor{
						Tag: "tag1",
						Inputs: []DataSpec{
							DataSpec{
								Type: "local",
							},
						},
					},
					CloudVendor{
						Tag: "tag2",
					},
				},
				RunTag: "tags",
			},
			want:  true,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				Kind:    tt.fields.Kind,
				Spec:    tt.fields.Spec,
				Vendors: tt.fields.Vendors,
			}
			got, got1 := job.HasLocals(tt.fields.RunTag)
			if got != tt.want {
				t.Errorf("Job.HasLocals() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Job.HasLocals() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
