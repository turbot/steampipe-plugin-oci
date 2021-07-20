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

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

variable "log_type" {
  type    = string
  default = "CUSTOM"
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_logging_log_group" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
}

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [oci_logging_log_group.named_test_resource]
  provisioner "local-exec" {
    command = "oci logging log create --display-name ${var.resource_name} --log-group-id ${oci_logging_log_group.named_test_resource.id} --log-type ${var.log_type}"
  }
  provisioner "local-exec" {
    command = "oci logging log list --log-group-id ${oci_logging_log_group.named_test_resource.id} --output json > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "display_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).data[0].id
}

output "log_group_id" {
  value = oci_logging_log_group.named_test_resource.id
}
