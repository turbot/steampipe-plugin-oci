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
}

resource "oci_cloud_guard_cloud_guard_configuration" "named_test_resource" {
    #Required
    compartment_id = var.tenancy_ocid
    reporting_region = var.region
    status = "ENABLED"
}

output "reporting_region" {
  value = var.region
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "self_manage_resources" {
  value = oci_cloud_guard_cloud_guard_configuration.named_test_resource.self_manage_resources
}
