variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
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
  tenancy_ocid = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_kms_vault" "named_test_resource" {
    compartment_id = var.tenancy_ocid
    display_name = var.resource_name
    vault_type = "DEFAULT"
}

resource "oci_kms_key" "named_test_resource" {
    depends_on  = [oci_kms_vault.named_test_resource]
    compartment_id = var.tenancy_ocid
    display_name = var.resource_name
    management_endpoint = oci_kms_vault.named_test_resource.management_endpoint
    key_shape {
        algorithm = "AES"
        length = 16
    }
}

resource "oci_kms_key_version" "named_test_resource" {
    depends_on  = [oci_kms_key.named_test_resource]
    key_id = oci_kms_key.named_test_resource.id
    management_endpoint = oci_kms_vault.named_test_resource.management_endpoint
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}

output "resource_id" {
  value = oci_kms_key_version.named_test_resource.id
}

output "key_id" {
  value = oci_kms_key.named_test_resource.id
}