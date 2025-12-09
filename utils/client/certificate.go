package client

// Certificate public and private certificate files
type Certificate struct {
	CertFile string `json:"public,omitempty"`
	KeyFile  string `json:"private,omitempty"`
}
