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

resource "oci_waf_web_app_firewall_policy" "test_policy" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
}

resource "oci_load_balancer_load_balancer" "test_lb" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  shape          = "flexible"
  subnet_ids     = [oci_core_subnet.test_subnet.id]

  shape_details {
    minimum_bandwidth_in_mbps = 10
    maximum_bandwidth_in_mbps = 100
  }
}

resource "oci_core_vcn" "test_vcn" {
  compartment_id = var.tenancy_ocid
  cidr_block     = "10.0.0.0/16"
  display_name   = var.resource_name
}

resource "oci_core_subnet" "test_subnet" {
  cidr_block     = "10.0.0.0/24"
  compartment_id = var.tenancy_ocid
  vcn_id         = oci_core_vcn.test_vcn.id
  display_name   = var.resource_name
}

resource "oci_waf_web_app_firewall" "test_resource" {
  compartment_id             = var.tenancy_ocid
  display_name               = var.resource_name
  backend_type               = "LOAD_BALANCER"
  load_balancer_id           = oci_load_balancer_load_balancer.test_lb.id
  web_app_firewall_policy_id = oci_waf_web_app_firewall_policy.test_policy.id
  freeform_tags              = { "Department" = "Finance" }
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
  value = oci_waf_web_app_firewall.test_resource.id
}

output "lifecycle_state" {
  value = oci_waf_web_app_firewall.test_resource.lifecycle_state
}

output "web_app_firewall_policy_id" {
  value = oci_waf_web_app_firewall_policy.test_policy.id
}
