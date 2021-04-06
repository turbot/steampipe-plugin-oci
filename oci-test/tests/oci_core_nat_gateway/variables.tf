variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "config_file_profile" {
  type        = string
  default     = "OCI"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

variable "policy_description" {
  type        = string
  default     = "Terraform testing resource"
  description = "The description you assign to the policy. Does not have to be unique, and it's changeable. "
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
