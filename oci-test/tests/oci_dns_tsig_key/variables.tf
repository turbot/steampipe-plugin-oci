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

variable "tsig_key_algorithm" {
  type        = string
  default     = "hmac-md5"
  description = "OCI tsig key algorithm used for the test."
}

variable "secret" {
  type        = string
  default     = "2020"
  description = "OCI secret used for the test."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_dns_tsig_key" "named_test_resource" {
  #Required
  algorithm      = var.tsig_key_algorithm
  compartment_id = var.tenancy_ocid
  name           = var.resource_name
  secret         = var.secret
}

output "resource_name" {
  value = oci_dns_tsig_key.named_test_resource.name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_dns_tsig_key.named_test_resource.id
}

