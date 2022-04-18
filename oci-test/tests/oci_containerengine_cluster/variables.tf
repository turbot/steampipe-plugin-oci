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

variable "cluster_kubernetes_version" {
  type        = string
  default     = "v1.22.5"
  description = "The version of Kubernetes to install into the cluster masters."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
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

resource "oci_containerengine_cluster" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  kubernetes_version = var.cluster_kubernetes_version
  name = var.resource_name
  vcn_id = oci_core_vcn.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_containerengine_cluster.named_test_resource.id
}

output "lifecycle_state" {
  value = oci_containerengine_cluster.named_test_resource.state
}

