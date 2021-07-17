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

variable "tag_namespace_description" {
  type        = string
  default     = "Terraform testing resource"
  description = "The description you assign to the tag namespace."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_identity_tag_namespace" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  description    = var.tag_namespace_description
  name           = var.resource_name
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "tag_namespace_description" {
  value = oci_identity_tag_namespace.named_test_resource.description
}

output "resource_id" {
  value = oci_identity_tag_namespace.named_test_resource.id
}
