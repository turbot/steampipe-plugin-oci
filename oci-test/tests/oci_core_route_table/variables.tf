variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  description = "OCI tenancy id."
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
}

resource "oci_core_route_table" "test_route_table" {
    depends_on  = [oci_core_internet_gateway.named_test_resource]
    compartment_id = var.tenancy_ocid
    vcn_id = oci_core_vcn.test_vcn.id
    display_name = var.resource_name
    freeform_tags = {"Department"= "Finance"}
    route_rules {
      network_entity_id = oci_core_internet_gateway.named_test_resource.id
      destination = "192.168.1.0/24"
    }
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_core_route_table.test_route_table.freeform_tags
}

output "resource_id" {
  value = oci_core_route_table.test_route_table.id
}

output "vcn_id" {
  value = oci_core_vcn.test_vcn.id
}
