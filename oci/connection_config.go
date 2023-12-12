package oci

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type ociConfig struct {
	Auth                  *string  `hcl:"auth"`
	ConfigPath            *string  `hcl:"config_path"`
	Fingerprint           *string  `hcl:"fingerprint"`
	PrivateKey            *string  `hcl:"private_key"`
	PrivateKeyPassword    *string  `hcl:"private_key_password"`
	PrivateKeyPath        *string  `hcl:"private_key_path"`
	Profile               *string  `hcl:"config_file_profile"`
	Regions               []string `hcl:"regions,optional"`
	TenancyOCID           *string  `hcl:"tenancy_ocid"`
	UserOCID              *string  `hcl:"user_ocid"`
	MaxErrorRetryAttempts *int     `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int     `hcl:"min_error_retry_delay"`
}

func ConfigInstance() interface{} {
	return &ociConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) ociConfig {
	if connection == nil || connection.Config == nil {
		return ociConfig{}
	}
	config, _ := connection.Config.(ociConfig)
	return config
}
