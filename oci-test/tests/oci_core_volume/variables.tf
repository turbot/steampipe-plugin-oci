variable "tenancy_ocid" {
  type        = string
  default     = ""
  description = "OCI tenancy id."
}

variable "config_file_profile" {
  type        = string
  default     = "OCI"
  description = "OCI credentials profile used for the test."
}

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI availability domain used for the test"
}

provider "oci" {
  tenancy_ocid = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_core_volume" "test_volume" {
  availability_domain = var.oci_ad
  compartment_id = var.tenancy_ocid
}

output "resource_id" {
  value = oci_core_volume.test_volume.id
}
