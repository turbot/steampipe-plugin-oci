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

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_identity_network_source" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  description    = var.resource_name
  name           = var.resource_name
  freeform_tags  = { "Name" = var.resource_name }
  provisioner "local-exec" {
    command = "sleep 120"
  }
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  depends_on = [oci_identity_network_source.named_test_resource]
  value = oci_identity_network_source.named_test_resource.freeform_tags
}

output "description" {
  value = var.resource_name
}

output "resource_id" {
  depends_on = [oci_identity_network_source.named_test_resource]
  value = oci_identity_network_source.named_test_resource.id
}
