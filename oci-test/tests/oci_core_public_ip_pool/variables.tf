variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
  description = "OCI tenancy id."
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
}

resource "oci_core_public_ip_pool" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  freeform_tags  = { "Name" = var.resource_name }
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}

output "resource_id" {
  value = oci_core_public_ip_pool.named_test_resource.id
}

output "resource_name" {
  value = oci_core_public_ip_pool.named_test_resource.display_name
}

output "lifecycle_state" {
  value = oci_core_public_ip_pool.named_test_resource.state
}

output "freeform_tags" {
  value = oci_core_public_ip_pool.named_test_resource.freeform_tags
}

output "cidr_blocks" {
  value = oci_core_public_ip_pool.named_test_resource.cidr_blocks
}
