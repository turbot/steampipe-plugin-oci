variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-HYDERABAD-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

locals {
  path = "${path.cwd}/output.json"
}
resource "oci_core_volume" "named_test_resource" {
  display_name        = var.resource_name
  availability_domain = "TvRS:AP-MUMBAI-1-AD-1"
  compartment_id      = var.tenancy_ocid
  block_volume_replicas {
    availability_domain = var.oci_ad
    display_name        = var.resource_name
  }
  block_volume_replicas_deletion = true
}

resource "null_resource" "named_test_resource" {
  depends_on = [oci_core_volume.named_test_resource]
  provisioner "local-exec" {
    command = "oci bv block-volume-replica list --availability-domain ${var.oci_ad} --compartment-id ${var.tenancy_ocid} --display-name ${var.resource_name} --output json > ${local.path}"
  }
  provisioner "local-exec" {
    command = "oci bv volume update --volume-id ${oci_core_volume.named_test_resource.id} --force --block-volume-replicas '[]'"
  }
}

data "local_file" "input" {
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
  value      = jsondecode(data.local_file.input.content).data[0].id
}

output "state" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).data[0].lifecycle-state
}
