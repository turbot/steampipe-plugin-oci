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
  default     = "Default"
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

resource "oci_core_vcn" "test_vcn" {
  compartment_id = var.tenancy_ocid
  cidr_block = "10.0.0.0/16"
}

resource "oci_core_internet_gateway" "named_test_resource" {
  depends_on  = [oci_core_vcn.test_vcn]
  compartment_id = var.tenancy_ocid
  vcn_id = oci_core_vcn.test_vcn.id
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
  value = oci_core_internet_gateway.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_core_internet_gateway.named_test_resource.id
}

output "display_name" {
  value = oci_core_internet_gateway.named_test_resource.display_name
}
