variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

data "oci_objectstorage_namespace" "test_namespace" {
  compartment_id = var.tenancy_ocid
}

resource "oci_objectstorage_bucket" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  name = var.resource_name
  namespace = data.oci_objectstorage_namespace.test_namespace.namespace
  freeform_tags = {"Department"= "Finance"}
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_objectstorage_bucket.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_objectstorage_bucket.named_test_resource.bucket_id
}

output "namespace" {
  value = data.oci_objectstorage_namespace.test_namespace.namespace
}
