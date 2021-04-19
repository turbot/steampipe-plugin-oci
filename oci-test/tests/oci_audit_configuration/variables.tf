variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "tenancy_ocid" {
  type        = string
  description = "OCID of your tenancy."
}

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_audit_configuration" "named_test_resource" {
    #Required
    compartment_id = var.tenancy_ocid
    retention_period_days = 365
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "retention_period_days" {
  value = oci_audit_configuration.named_test_resource.retention_period_days
}

