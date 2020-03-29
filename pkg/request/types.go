package request

// LoginResponse defines the body from login request
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

// CreateRequest defines the create request
type CreateRequest struct {
	UserName string `json:"user_name,omitempty"`
	YAML     string `json:"yaml,omitempty"`
	Name     string `json:"name,omitempty"`
	Image    string `json:"image,omitempty"`
}
