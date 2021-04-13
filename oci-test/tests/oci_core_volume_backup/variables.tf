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

resource "oci_core_volume" "test_volume" {
  availability_domain = var.oci_ad
  compartment_id = var.tenancy_ocid
}

resource "oci_core_volume_backup" "test_volume_backup" {
  depends_on  = [oci_core_volume.test_volume]
  volume_id = oci_core_volume.test_volume.id
  display_name = var.resource_name
  freeform_tags = {"Department"= "Finance"}
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_core_volume_backup.test_volume_backup.freeform_tags
}

output "resource_id" {
  value = oci_core_volume_backup.test_volume_backup.id
}

output "volume_id" {
  value = oci_core_volume.test_volume.id
}

