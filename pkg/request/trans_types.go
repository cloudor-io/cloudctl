package request

// TransInfo defines a transaction
type TransInfo struct {
	// job or pay
	Type         string  `json:"type,omitempty"`
	JobID        string  `json:"job_id,omitempty"`
	CreditBefore float64 `json:"credit_before,omitempty"`
	CreditAfter  float64 `json:"credit_after,omitempty"`
}

// TransSchema defines the json returned from transaction list
type TransSchema struct {
	UserName  string    `json:"user_name,omitempty"`
	ID        string    `json:"id,omitempty"`
	TimeStamp int64     `json:"timestamp,omitempty"`
	Amount    float64   `json:"amount"`
	Unit      string    `json:"unit,omitempty"`
	Info      TransInfo `json:"info,omitempty"`
}
