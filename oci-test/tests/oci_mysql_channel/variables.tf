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

resource "oci_core_vcn" "named_test_resource" {
    compartment_id = var.tenancy_ocid
    cidr_block = "10.0.0.0/16"
}
resource "oci_core_subnet" "named_test_resource" {
    cidr_block = "10.0.0.0/16"
    compartment_id = var.tenancy_ocid
    vcn_id = oci_core_vcn.named_test_resource.id
}

resource "oci_mysql_mysql_db_system" "named_test_resource" {
    #Required
    admin_password ="Admin@1234"
    admin_username = "admin"
    availability_domain = "TvRS:AP-MUMBAI-1-AD-1"
    compartment_id = var.tenancy_ocid
    shape_name = "VM.Standard.E2.1"
    subnet_id = oci_core_subnet.named_test_resource.id
    display_name = var.resource_name
    data_storage_size_in_gb = "50"
}

resource "oci_mysql_channel" "named_test_resource" {
    #Required
    source {
        #Required
        hostname = var.resource_name
        password = "Test@1234"
        source_type = "MYSQL"
        ssl_mode = "REQUIRED"
        username = var.resource_name
    }

    target {
        #Required
        db_system_id = oci_mysql_mysql_db_system.named_test_resource.id
        target_type = "DBSYSTEM"
    }

    #Optional
    description = var.resource_name
    display_name = var.resource_name
    freeform_tags = {"name"= var.resource_name}
    is_enabled = true
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
  value = oci_mysql_channel.named_test_resource.id
}

output "state" {
  value = oci_mysql_channel.named_test_resource.state
}

output "target" {
  value = oci_mysql_channel.named_test_resource.target
}

output "time_created" {
  value = oci_mysql_channel.named_test_resource.time_created
}

output "description" {
  value = oci_mysql_channel.named_test_resource.description
}