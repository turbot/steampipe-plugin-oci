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

resource "oci_core_service_gateway" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  services {
    #Required
    service_id = data.oci_core_services.named_test_resource.services.0.id
  }
  vcn_id = oci_core_vcn.named_test_resource.id

  #Optional
  display_name   = var.resource_name
  freeform_tags  = { "Name" = var.resource_name }
  route_table_id = oci_core_route_table.named_test_resource.id
}

resource "oci_core_vcn" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name

  #Optional
  cidr_block = "10.0.0.0/24"
}

resource "oci_core_route_table" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  vcn_id         = oci_core_vcn.named_test_resource.id
}

data "oci_core_services" "named_test_resource" {
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
  value = oci_core_service_gateway.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_core_service_gateway.named_test_resource.id
}

output "route_table_id" {
  value = oci_core_service_gateway.named_test_resource.route_table_id
}

output "state" {
  value = oci_core_service_gateway.named_test_resource.state
}

output "vcn_id" {
  value = oci_core_service_gateway.named_test_resource.vcn_id
}

output "block_traffic" {
  value = oci_core_service_gateway.named_test_resource.block_traffic
}