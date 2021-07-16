variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
  description = "OCID of your tenancy."
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

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_apigateway_api" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid

  #Optional
  display_name = var.resource_name
}

output "reporting_region" {
  value = var.region
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = oci_apigateway_api.named_test_resource.id
}
