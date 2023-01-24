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
  freeform_tags       = { "Name" = var.resource_name }
}

resource "oci_core_volume_group" "test_volume_group" {
    #Required
    availability_domain = var.oci_ad
    compartment_id = var.tenancy_ocid
    source_details {
        #Required
        type = "volumeIds"
        volume_ids = [oci_core_volume.test_volume.id]
    }

    #Optional
    # defined_tags = {"Name"= var.resource_name}
    display_name = var.resource_name
    freeform_tags = {"Department"= var.resource_name}
    # volume_group_replicas {
    #     #Required
    #     availability_domain = var.oci_ad

    #     #Optional
    #     display_name = var.resource_name
    # }
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_core_volume_group.test_volume_group.id
}

output "display_name" {
  value = var.resource_name
}
