variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
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
  depends_on     = [oci_core_vcn.named_test_resource]
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
  vcn_id         = oci_core_vcn.named_test_resource.id
  freeform_tags  = { "Name" = var.resource_name }
}

resource "oci_file_storage_mount_target" "named_test_resource" {
  depends_on          = [oci_core_subnet.named_test_resource]
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
  subnet_id           = oci_core_subnet.named_test_resource.id
  display_name        = var.resource_name
  freeform_tags       = { "Name" = var.resource_name }
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
  value = oci_file_storage_mount_target.named_test_resource.id
}

output "availability_domain" {
  value = oci_file_storage_mount_target.named_test_resource.availability_domain
}

output "subnet_id" {
  value = oci_file_storage_mount_target.named_test_resource.subnet_id
}

output "state" {
  value = oci_file_storage_mount_target.named_test_resource.state
}

output "freeform_tags" {
  value = oci_file_storage_mount_target.named_test_resource.freeform_tags
}
