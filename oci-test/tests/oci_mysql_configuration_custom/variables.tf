variable "resource_name" {
  type        = string
  default     = "turbot-test-20210818-config-custom"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
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
    command = "oci mysql configuration create --compartment-id ${var.tenancy_ocid} --shape-name ${var.shape} --output json > ${local.path}"
  }
}

data "local_file" "configuration" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}

output "resource_name" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data.display-name
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data.id
}

output "type" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data.type
}

output "shape_name" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data.shape-name
}

output "lifecycle_state" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.configuration.content).data.lifecycle-state
}
