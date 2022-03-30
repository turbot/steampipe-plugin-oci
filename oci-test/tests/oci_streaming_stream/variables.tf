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
  default     = ""
  description = "OCID of your tenancy."
}

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

variable "stream_partitions" {
  type        = string
  default     = "1"
  description = "The number of partitions in the stream."
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

resource "oci_streaming_stream_pool" "test_stream_pool" {
    #Required
    compartment_id = var.tenancy_ocid
    name = var.resource_name
}

resource "oci_streaming_stream" "test_stream" {
    #Required
    name = var.resource_name
    partitions = var.stream_partitions

    #Optional
    freeform_tags = {"Department"= "Finance"}
    stream_pool_id = oci_streaming_stream_pool.test_stream_pool.id
}


output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_streaming_stream.test_stream.id
}

output "state" {
  value = oci_streaming_stream.test_stream.state
}