variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = ""
  description = "OCID of your tenancy."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

locals {
  path       = "${path.cwd}/identity_tenancy.json"
  periodPath = "${path.cwd}/retention_period.json"
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "oci iam tenancy get --tenancy-id ${var.tenancy_ocid} > ${local.path}"
  }
}

resource "null_resource" "retention_period" {
  provisioner "local-exec" {
    command = "oci audit config get --compartment-id ${var.tenancy_ocid} > ${local.periodPath}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

data "local_file" "periodInput" {
  depends_on = [null_resource.retention_period]
  filename   = local.periodPath
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "data", "data"), "id", "id")
}

output "resourceName" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "data", "data"), "name", "name")
}

output "description" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "data", "data"), "description", "description")
}

output "retention_period_days" {
  depends_on = [null_resource.retention_period]
  value      = lookup(lookup(jsondecode(data.local_file.periodInput.content), "data", "data"), "retention-period-days", "retention-period-days")
}
