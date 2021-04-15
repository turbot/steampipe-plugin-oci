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
  default     = " DEFAULT"
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

resource "oci_core_vcn" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name = var.resource_name
  cidr_block = "10.0.0.0/16"
}

resource "oci_core_subnet" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name = var.resource_name
  cidr_block = "10.0.0.0/16"
  vcn_id = oci_core_vcn.named_test_resource.id
}

resource "oci_functions_application" "test_application" {
  compartment_id = var.tenancy_ocid
  display_name = var.resource_name
  subnet_ids = [oci_core_subnet.named_test_resource.id]
  freeform_tags = {"Department"= "Finance"}
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}

output "freeform_tags" {
  value = oci_functions_application.test_application.freeform_tags
}

output "resource_id" {
  value = oci_functions_application.test_application.id
}

output "display_name" {
  value = oci_functions_application.test_application.display_name
}

output "subnet_id" {
  value = oci_core_subnet.named_test_resource.id
}
