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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_core_volume" "test_volume" {
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
  freeform_tags       = { "Name" = var.resource_name }
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_core_volume.test_volume.freeform_tags
}

output "resource_id" {
  value = oci_core_volume.test_volume.id
}
