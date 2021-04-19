variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
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

resource "oci_logging_log_group" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name = var.resource_name
}

resource "oci_logging_log" "named_test_resource" {
  log_group_id = oci_logging_log_group.named_test_resource.id
  display_name = var.resource_name
  log_type = "Custom log"
}

output "display_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}

output "resource_id" {
  value = oci_logging_log.named_test_resource.id
}

output "log_group_id" {
  value = oci_logging_log.named_test_resource.log_group_id
}