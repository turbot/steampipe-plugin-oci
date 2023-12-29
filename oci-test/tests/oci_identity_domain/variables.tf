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

variable "domain_description" {
  type        = string
  default     = "Terraform testing resource"
  description = "The description you assign to the domain. Does not have to be unique, and it's changeable. "
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}


resource "oci_identity_domain" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  description = var.domain_description
  display_name = var.resource_name
  home_region = var.region
  license_type = "Free"

  #Optional
  freeform_tags = { "Name" = var.resource_name }
  provisioner "local-exec" {
    command = "sleep 60"
  }
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  depends_on = [oci_identity_domain.named_test_resource]
  value = oci_identity_domain.named_test_resource.freeform_tags
}

output "resource_id" {
  depends_on = [oci_identity_domain.named_test_resource]
  value = oci_identity_domain.named_test_resource.id
}

output "license_type" {
  depends_on = [oci_identity_domain.named_test_resource]
  value = oci_identity_domain.named_test_resource.license_type
}

output "domain_description" {
  depends_on = [oci_identity_domain.named_test_resource]
  value = oci_identity_domain.named_test_resource.description
}

output "time_created" {
  depends_on = [oci_identity_domain.named_test_resource]
  value = oci_identity_domain.named_test_resource.time_created
}
