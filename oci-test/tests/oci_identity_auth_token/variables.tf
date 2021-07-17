variable "user_name" {
  type        = string
  default     = "steampipetest20200125user"
  description = "Name of the resource used throughout the test."
}

variable "resource_name" {
  type        = string
  default     = "steampipetest20200125e"
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

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_identity_user" "test_user" {
  compartment_id = var.tenancy_ocid
  description    = var.user_name
  name           = var.user_name
}

resource "oci_identity_auth_token" "test_auth_token" {
  description = var.resource_name
  user_id     = oci_identity_user.test_user.id
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_name" {
  value = var.resource_name
}

output "user_id" {
  value = oci_identity_auth_token.test_auth_token.user_id
}

output "resource_id" {
  value = oci_identity_auth_token.test_auth_token.id
}
