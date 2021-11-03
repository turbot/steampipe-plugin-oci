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

variable "dynamic_group_matching_rule" {
  type        = string
  default     = "Any {instance.id = 'ocid1.instance.oc1.iad..exampleuniqueid1', instance.compartment.id = 'ocid1.compartment.oc1..exampleuniqueid2'}"
  description = "The matching rule you assign to the dynamic group."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_identity_dynamic_group" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  description    = var.resource_name
  matching_rule  = var.dynamic_group_matching_rule
  name           = var.resource_name
  freeform_tags  = { "Department" = "Finance" }
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

output "description" {
  depends_on = [oci_identity_dynamic_group.named_test_resource]
  value = oci_identity_dynamic_group.named_test_resource.description
}

output "freeform_tags" {
  depends_on = [oci_identity_dynamic_group.named_test_resource]
  value = oci_identity_dynamic_group.named_test_resource.freeform_tags
}

output "resource_id" {
  depends_on = [oci_identity_dynamic_group.named_test_resource]
  value = oci_identity_dynamic_group.named_test_resource.id
}

