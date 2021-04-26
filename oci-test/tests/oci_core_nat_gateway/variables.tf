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

resource "oci_core_nat_gateway" "named_test_resource" {
  depends_on = [
    oci_core_vcn.named_test_resource
  ]
  #Required
  compartment_id = var.tenancy_ocid
  vcn_id         = oci_core_vcn.named_test_resource.id

  #Optional
  block_traffic = "true"
  display_name  = var.resource_name
  freeform_tags = { "Name" = var.resource_name }
}

resource "oci_core_vcn" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  display_name = var.resource_name

  #Optional
  cidr_block = "10.0.0.0/16"
}

output "vcn_id" {
  value = oci_core_vcn.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_core_nat_gateway.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_core_nat_gateway.named_test_resource.id
}

output "time_created" {
  value = oci_core_nat_gateway.named_test_resource.time_created
}

output "display_name" {
  value = oci_core_nat_gateway.named_test_resource.display_name
}
