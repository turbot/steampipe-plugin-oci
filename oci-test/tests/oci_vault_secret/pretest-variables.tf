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

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_kms_vault" "named_test_resource" {
  compartment_id   = var.tenancy_ocid
  display_name     = var.resource_name
  vault_type       = "DEFAULT"
}

resource "oci_kms_key" "named_test_resource" {
  depends_on          = [oci_kms_vault.named_test_resource]
  compartment_id      = var.tenancy_ocid
  display_name        = var.resource_name
  management_endpoint = oci_kms_vault.named_test_resource.management_endpoint
  key_shape {
    algorithm = "AES"
    length    = 16
  }
}

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [oci_kms_key.named_test_resource]
  provisioner "local-exec" {
    command = "oci vault secret create-base64 --compartment-id ${var.tenancy_ocid} --secret-name ${var.resource_name} --vault-id ${oci_kms_vault.named_test_resource.id} --key-id ${oci_kms_key.named_test_resource.id} --secret-content-content VGhpcyBpcyB0ZXN0IHNlY3JldA== --output json > ${local.path}"
  }
  provisioner "local-exec" {
    command = "sleep 60"
  }
}

data "local_file" "secret" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.secret.content).data.id
}

output "lifecycle_state" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.secret.content).data.lifecycle-state
}
