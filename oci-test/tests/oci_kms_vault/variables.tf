variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
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
}

resource "oci_kms_vault" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  vault_type     = "DEFAULT"
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

output "resource_id" {
  value = oci_kms_vault.named_test_resource.id
}

output "vault_type" {
  value = oci_kms_vault.named_test_resource.vault_type
}

output "crypto_endpoint" {
  value = oci_kms_vault.named_test_resource.crypto_endpoint
}

output "management_endpoint" {
  value = oci_kms_vault.named_test_resource.management_endpoint
}
