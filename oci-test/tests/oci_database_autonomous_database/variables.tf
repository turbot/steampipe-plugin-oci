variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "tenancy_ocid" {
  type        = string
  description = "OCID of your tenancy."
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

resource "oci_database_autonomous_database" "named_test_resource" {
    #Required
    compartment_id = var.tenancy_ocid
    cpu_core_count = 1
    display_name = var.resource_name
    admin_password = "Test@12345678"
    data_storage_size_in_tbs = 1

    #The Autonomous Database name cannot be longer than 14 characters.
    db_name = "dbTest"
}

output "resource_name" {
  value = var.resource_name
}

output "db_name" {
  value = oci_database_autonomous_database.named_test_resource.db_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_database_autonomous_database.named_test_resource.id
}
