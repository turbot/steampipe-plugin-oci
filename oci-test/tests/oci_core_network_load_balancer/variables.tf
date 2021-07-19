variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
  description = "OCI credentials profile used for the test. Default is to use the default profile."
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

resource "oci_core_vcn" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
}

resource "oci_core_subnet" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
  vcn_id         = oci_core_vcn.named_test_resource.id
}

resource "oci_network_load_balancer_network_load_balancer" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  subnet_id      = oci_core_subnet.named_test_resource.id
  freeform_tags  = { "Name" = var.resource_name }
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
  value = oci_network_load_balancer_network_load_balancer.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_network_load_balancer_network_load_balancer.named_test_resource.id
}

output "display_name" {
  value = oci_network_load_balancer_network_load_balancer.named_test_resource.display_name
}

output "subnet_id" {
  value = oci_network_load_balancer_network_load_balancer.named_test_resource.subnet_id
}

output "is_preserve_source_destination" {
  value = oci_network_load_balancer_network_load_balancer.named_test_resource.is_preserve_source_destination
}

output "is_private" {
  value = oci_network_load_balancer_network_load_balancer.named_test_resource.is_private
}

output "state" {
  value = oci_network_load_balancer_network_load_balancer.named_test_resource.state
}
