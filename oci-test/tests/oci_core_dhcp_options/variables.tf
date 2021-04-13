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

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_core_vcn" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  cidr_block     = "10.0.0.0/16"
}

resource "oci_core_dhcp_options" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  options {
    type               = "DomainNameServer"
    server_type        = "CustomDnsServer"
    custom_dns_servers = ["192.168.0.2", "192.168.0.11", "192.168.0.19"]
  }

  options {
    type                = "SearchDomain"
    search_domain_names = ["test.com"]
  }

  vcn_id = oci_core_vcn.named_test_resource.id

  #Optional
  display_name  = var.resource_name
  freeform_tags = { "Department" = "Finance" }
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

output "vcn_id" {
  value = oci_core_dhcp_options.named_test_resource.vcn_id
}

output "freeform_tags" {
  value = oci_core_dhcp_options.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_core_dhcp_options.named_test_resource.id
}

output "display_name" {
  value = oci_core_dhcp_options.named_test_resource.display_name
}
