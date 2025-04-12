package models

// runtime config struct
type WorkmanConfig struct {
	CONFIG_PATH string
}

// InstanceDetails represents the structure of each entry in the JSON file
type InstanceDetails struct {
	ID           string `json:"id"`
	PEM          string `json:"pem"`
	AwsProfile   string `json:"aws_profile"`              // AWS profile for session creation
	InstanceUser string `json:"instance_user"`            // SSH user for the instance
	UsePrivateIP bool   `json:"use_private_ip,omitempty"` // Optional: Use private IP for SSH
	LastAccessed string `json:"last_accessed,omitempty"`  // Timestamp of last access
}
