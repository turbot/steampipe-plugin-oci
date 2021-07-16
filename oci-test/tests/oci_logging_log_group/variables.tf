variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
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

resource "oci_logging_log_group" "test_log_group" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  description    = var.resource_name
  freeform_tags  = { "Name" = var.resource_name }
}

output "display_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_logging_log_group.test_log_group.freeform_tags
}

output "region" {
  value = var.region
}

output "resource_id" {
  value = oci_logging_log_group.test_log_group.id
}
