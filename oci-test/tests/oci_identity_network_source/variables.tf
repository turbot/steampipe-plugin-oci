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

resource "oci_identity_group" "test_group" {
  #Required
  compartment_id = var.tenancy_ocid
  description    = var.resource_name
  name           = var.resource_name
}


resource "oci_identity_network_source" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  description = var.resource_name
  name = var.resource_name
  freeform_tags = {"Name"= var.resource_name}
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_identity_network_source.named_test_resource.freeform_tags
}

output "description" {
  value = var.resource_name
}

output "resource_id" {
  value = oci_identity_network_source.named_test_resource.id
}
