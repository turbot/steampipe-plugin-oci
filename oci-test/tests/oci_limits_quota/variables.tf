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

variable "quota_statements" {
  type = list
  default = ["Zero notifications quotas in tenancy"]
  description = "An array of one or more quota statements written in the declarative quota statement language."
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

resource "oci_limits_quota" "named_test_resource" {
    #Required
    compartment_id = var.tenancy_ocid
    description = var.resource_name
    name = var.resource_name
    statements = var.quota_statements

    #Optional
    freeform_tags = {"Department"= "Finance"}
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
  value = oci_limits_quota.named_test_resource.id
}

output "state" {
  value = oci_limits_quota.named_test_resource.state
}