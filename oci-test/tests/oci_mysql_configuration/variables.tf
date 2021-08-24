variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
  description = "OCI tenancy id."
}

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region"
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

variable "shape" {
  type        = string
  default     = "VM.Standard.E2.1"
  description = "Oracle My SQL Shape."
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
    command = "oci mysql configuration list --compartment-id ${var.tenancy_ocid} --output json > ${local.path}"
  }
}

data "local_file" "configuration" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_name" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data[0].display-name
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data[0].id
}

output "type" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data[0].type
}

output "shape_name" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data[0].shape-name
}

output "lifecycle_state" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data[0].lifecycle-state
}
