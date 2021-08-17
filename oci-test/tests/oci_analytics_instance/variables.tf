variable "resource_name" {
  type        = string
  default     = "steampipetest20210817"
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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_analytics_analytics_instance" "test_analytics_instance" {
    #Required
    capacity {
        #Required
        capacity_type = "OLPU_COUNT"
        capacity_value = 1
    }
    compartment_id = var.tenancy_ocid
    feature_set = "SELF_SERVICE_ANALYTICS"
    idcs_access_token = var.resource_name
    license_type = "BRING_YOUR_OWN_LICENSE"
    name = var.resource_name
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_analytics_analytics_instance.test_analytics_instance.id
}

