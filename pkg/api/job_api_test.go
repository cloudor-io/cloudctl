package api

import (
	"testing"
)

func TestJob_FindRunningVendorIndexByTag(t *testing.T) {

	type fields struct {
		Kind    string
		UUID    string
		Name    string
		RunTag  string
		Spec    RunSpec
		Vendors []CloudVendor
	}
	exampleField := fields{
		RunTag: "tag",
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
		want   int
	}{
		// TODO: Add test cases.
		{
			name:   "no_vendor_should_return_negative",
			fields: fields{},
			want:   -1,
		},
		{
			name:   "no_tag_should_return_index_0",
			fields: exampleField,
			want:   0,
		},
		{
			name: "tag2_should_return_index_1",
			fields: fields{
				RunTag: "tag2",
				Vendors: []CloudVendor{
					CloudVendor{
						Tag: "tag1",
					},
					CloudVendor{
						Tag: "tag2",
					},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				Kind:    tt.fields.Kind,
				RunTag:  tt.fields.RunTag,
				Spec:    tt.fields.Spec,
				Vendors: tt.fields.Vendors,
			}
			if got := job.FindRunningVendorIndexByTag(); got != tt.want {
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
		RunTag  string
		Spec    RunSpec
		Vendors []CloudVendor
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
				RunTag: "tag1",
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
			},
			want:  true,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				Kind:    tt.fields.Kind,
				RunTag:  tt.fields.RunTag,
				Spec:    tt.fields.Spec,
				Vendors: tt.fields.Vendors,
			}
			got, got1 := job.HasLocals()
			if got != tt.want {
				t.Errorf("Job.HasLocals() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Job.HasLocals() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
