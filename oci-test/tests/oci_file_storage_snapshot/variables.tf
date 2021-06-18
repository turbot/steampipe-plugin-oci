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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_file_storage_file_system" "named_test_resource" {
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
  display_name        = var.resource_name
  freeform_tags       = { "Name" = var.resource_name }
}

resource "oci_file_storage_snapshot" "named_test_resource" {
  file_system_id = oci_file_storage_file_system.named_test_resource.id
  name           = var.resource_name
  freeform_tags  = { "Name" = var.resource_name }
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_file_storage_snapshot.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_file_storage_snapshot.named_test_resource.id
}

output "file_system_id" {
  value = oci_file_storage_snapshot.named_test_resource.file_system_id
}

output "lifecycle_state" {
  value = oci_file_storage_snapshot.named_test_resource.state
}