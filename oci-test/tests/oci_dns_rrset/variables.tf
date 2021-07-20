variable "resource_name" {
  type        = string
  default     = "steampipetest"
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

variable "oci_ad" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_dns_zone" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  name           = "steampipetest.com"
  zone_type      = "PRIMARY"
  scope          = "GLOBAL"
  freeform_tags  = { "Name" = "steampipetest" }
}

resource "oci_dns_rrset" "named_test_resource" {
  #Required
  zone_name_or_id = oci_dns_zone.named_test_resource.name
  domain          = "test"
  rtype           = "NS"

  #Optional
  compartment_id = var.tenancy_ocid
}

output "domain" {
  value = oci_dns_zone.named_test_resource.name
}

output "rtype" {
  value = oci_dns_rrset.named_test_resource.rtype
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}
