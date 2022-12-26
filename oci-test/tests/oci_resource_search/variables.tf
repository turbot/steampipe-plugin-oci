variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
  description = "OCID of your tenancy."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

variable "query" {
  type        = string
  default     = "query user resources"
  description = "The query based on which search was performed."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "oci search resource structured-search --query-text '${var.query}' --output json > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_name" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).data.items[0].display-name
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).data.items[0].identifier
}

output "query" {
  depends_on = [null_resource.named_test_resource]
  value      = var.query
}

output "region" {
  value = var.region
}
