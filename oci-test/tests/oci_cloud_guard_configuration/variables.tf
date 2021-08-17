variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

# we are not creating any resource so we have to pass the reporting region where resorce will be available
variable "region" {
  type        = string
  default     = "ap-hyderabad-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "oci cloud-guard configuration get --compartment-id ${var.tenancy_ocid} --output json > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "reporting_region" {
  value = var.region
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "self_manage_resources" {
  depends_on = [
     null_resource.named_test_resource
  ]
  value = jsondecode(data.local_file.input.content).data.self-manage-resources
}

output "status" {
  depends_on = [
     null_resource.named_test_resource
  ]
  value = jsondecode(data.local_file.input.content).data.status
}
