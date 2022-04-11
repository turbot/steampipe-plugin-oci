package oci

import (
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/schema"
)

type ociConfig struct {
	Auth                  *string  `cty:"auth"`
	ConfigPath            *string  `cty:"config_path"`
	Fingerprint           *string  `cty:"fingerprint"`
	PrivateKey            *string  `cty:"private_key"`
	PrivateKeyPassword    *string  `cty:"private_key_password"`
	PrivateKeyPath        *string  `cty:"private_key_path"`
	Profile               *string  `cty:"config_file_profile"`
	Regions               []string `cty:"regions"`
	TenancyOCID           *string  `cty:"tenancy_ocid"`
	UserOCID              *string  `cty:"user_ocid"`
	MaxErrorRetryAttempts *int     `cty:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int     `cty:"min_error_retry_delay"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"regions": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"auth": {
		Type: schema.TypeString,
	},
	"tenancy_ocid": {
		Type: schema.TypeString,
	},
	"config_file_profile": {
		Type: schema.TypeString,
	},
	"config_path": {
		Type: schema.TypeString,
	},
	"user_ocid": {
		Type: schema.TypeString,
	},
	"fingerprint": {
		Type: schema.TypeString,
	},
	"private_key": {
		Type: schema.TypeString,
	},
	"private_key_path": {
		Type: schema.TypeString,
	},
	"private_key_password": {
		Type: schema.TypeString,
	},
	"max_error_retry_attempts": {
		Type: schema.TypeInt,
	},
	"min_error_retry_delay": {
		Type: schema.TypeInt,
	},
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
