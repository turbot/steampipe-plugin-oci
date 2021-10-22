variable "resource_name" {
  type        = string
  default     = "turbot-test-20210819-database-software-image"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
  description = "OCID of your tenancy."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "oci_ad" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_database_database_software_image" "named_test_resource" {
  compartment_id   = var.tenancy_ocid
  display_name     = var.resource_name
  database_version = "19.0.0.0"
  freeform_tags    = { "Name" = var.resource_name }
  image_type       = "DATABASE_IMAGE"
  patch_set        = "19.11.0.0"
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "image_type" {
  value = oci_database_database_software_image.named_test_resource.image_type
}

output "resource_id" {
  value = oci_database_database_software_image.named_test_resource.id
}

output "database_version" {
  value = oci_database_database_software_image.named_test_resource.database_version
}

output "patch_set" {
  value = oci_database_database_software_image.named_test_resource.patch_set
}
