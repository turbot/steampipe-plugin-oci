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

variable "heat_wave_cluster_cluster_size" {
  type        = number
  default     = 2
  description = "The number of analytics-processing compute instances, of the specified shape, in the HeatWave cluster."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
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

resource "oci_core_vcn" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  cidr_block     = "10.0.0.0/16"
}

resource "oci_core_subnet" "named_test_resource" {
  cidr_block     = "10.0.0.0/16"
  compartment_id = var.tenancy_ocid
  vcn_id         = oci_core_vcn.named_test_resource.id
}

resource "oci_mysql_mysql_db_system" "named_test_resource" {
  #Required
  admin_password          = "Admin@1234"
  admin_username          = "admin"
  availability_domain     = "TvRS:AP-MUMBAI-1-AD-1"
  compartment_id          = var.tenancy_ocid
  shape_name              = "MySQL.HeatWave.VM.Standard.E3"
  subnet_id               = oci_core_subnet.named_test_resource.id
  display_name            = var.resource_name
  data_storage_size_in_gb = 50
}

resource "oci_mysql_heat_wave_cluster" "named_test_resource" {
  #Required
  db_system_id = oci_mysql_mysql_db_system.named_test_resource.id
  cluster_size = var.heat_wave_cluster_cluster_size
  shape_name = "MySQL.HeatWave.VM.Standard.E3"
}

output "resource_id" {
  value = oci_mysql_heat_wave_cluster.named_test_resource.db_system_id
}

output "lifecycle_state" {
  value = oci_mysql_heat_wave_cluster.named_test_resource.state
}

output "shape_name" {
  value = oci_mysql_heat_wave_cluster.named_test_resource.shape_name
}