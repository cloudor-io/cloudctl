package request

// Vendors defines a multi-level map of vendor instances
// vendor_name -> region -> instance_type -> Instance
type Vendors map[string]map[string]map[string]Instance

// Instance specifies instance features
type Instance struct {
	Price     float64 `json:"price,omitempty"`
	Unit      string  `json:"unit,omitempty"`
	CPU       int32   `json:"cpu,omitempty"`
	GPU       string  `json:"gpu,omitempty"`
	MemoryGiB int32   `json:"memory_gib,omitempty"`
}
