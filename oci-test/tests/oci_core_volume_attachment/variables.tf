variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
  description = "OCI tenancy id."
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
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_core_volume" "test_volume" {
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
}

resource "oci_core_instance" "test_instance" {
  #Required
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
  shape               = "VM.Standard.E4.Flex"
}

resource "oci_core_volume_attachment" "test_volume" {
  #Required
  attachment_type = "iscsi"
  instance_id     = oci_core_instance.test_instance.id
  volume_id       = oci_core_volume.test_volume.id
  display_name    = var.resource_name
  is_read_only    = false
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_core_volume_attachment.test_volume.id
}
