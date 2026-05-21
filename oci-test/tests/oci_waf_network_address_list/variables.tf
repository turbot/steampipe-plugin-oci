variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
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

resource "oci_waf_network_address_list" "test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  type           = "ADDRESSES"
  addresses      = ["192.168.0.0/24"]
  freeform_tags  = { "Department" = "Finance" }
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

output "resource_id" {
  value = oci_waf_network_address_list.test_resource.id
}

output "lifecycle_state" {
  value = oci_waf_network_address_list.test_resource.lifecycle_state
}
