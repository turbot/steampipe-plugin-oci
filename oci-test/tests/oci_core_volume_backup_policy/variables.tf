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
  default     = ""
  description = "OCID of your tenancy."
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

resource "oci_core_volume_backup_policy" "named_test_resource" {

  #Optional
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  freeform_tags  = { "Name" = var.resource_name }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = oci_core_volume_backup_policy.named_test_resource.id
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}
