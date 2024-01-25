variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
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

resource "oci_identity_user" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  description    = var.resource_name
  name           = var.resource_name
}

resource "oci_identity_db_credential" "named_test_resource" {
    #Required
    description = "testing"
    password = "TurbotKolkata@123"
    user_id = oci_identity_user.named_test_resource.id
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "user_id" {
  value = oci_identity_db_credential.named_test_resource.user_id
}

output "credential_id" {
  value = split("dbCredentials/",oci_identity_db_credential.named_test_resource.id)[1]
}

output "lifecycle_state" {
  value = oci_identity_db_credential.named_test_resource.state
}

output "description" {
  value = oci_identity_db_credential.named_test_resource.description
}

