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

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_core_vcn" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
}

resource "oci_core_security_list" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  vcn_id         = oci_core_vcn.named_test_resource.id
  display_name   = var.resource_name
  freeform_tags  = { "Name" = var.resource_name }
}

output "resource_name" {
  value = var.resource_name
}

output "region" {
  value = var.region
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_core_security_list.named_test_resource.id
}

output "vcn_id" {
  value = oci_core_security_list.named_test_resource.vcn_id
}
