package impl

import (
	"encoding/json"
	"fmt"

	"github.com/cloudor-io/cloudctl/pkg/request"
)

// GetVendors returns supported vendors by cloudor
func GetVendors(userName, token *string) (*request.Vendors, error) {
	resp, err := request.GetCloudor(userName, token, "/vendor")
	if err != nil {
		fmt.Printf("getting vendors failed for user %s: %v", *userName, err)
		return nil, err
	}

	vendors := request.Vendors{}
	err = json.Unmarshal(resp, &vendors)
	if err != nil {
		fmt.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	return &vendors, nil
}

// Instance specifies instance features
type Instance struct {
	Vendor       string           `json:"vendor,omitempty"`
	Region       string           `json:"region,omitempty"`
	InstanceType string           `json:"instance_type,omitempty"`
	Instance     request.Instance `json:"instance,omitempty"`
}

func FilterVendors(vendors *request.Vendors, selVendor string, selRegion string) []Instance {
	instances := []Instance{}
	for vendorName, vendor := range *vendors {
		if selVendor != "" {
			if vendorName != selVendor {
				continue
			}
		}
		for regionName, region := range vendor {
			if selRegion != "" {
				if regionName != selRegion {
					continue
				}
			}
			for instanceType, instance := range region {
				instances = append(instances, Instance{
					Vendor:       vendorName,
					Region:       regionName,
					InstanceType: instanceType,
					Instance:     instance,
				})
			}
		}
	}
	return instances
}
