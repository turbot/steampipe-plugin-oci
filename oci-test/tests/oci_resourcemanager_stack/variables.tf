variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
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

variable "terraform_version" {
  type        = string
  default     = "0.12.x"
  description = "The version of Terraform to use with the stack."
}

variable "template_category_id" {
  type        = string
  default     = "0"
  description = "The category in which the template belongs to"
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

locals {
  template_path = "${path.cwd}/template-output.json"
  resource_path = "${path.cwd}/resource-output.json"
}

resource "null_resource" "named_test_template" {
  provisioner "local-exec" {
    command = "oci resource-manager template list --compartment-id ${var.tenancy_ocid} --template-category-id ${var.template_category_id} --output json > ${local.template_path}"
  }
}

data "local_file" "input_template" {
  depends_on = [null_resource.named_test_template]
  filename   = local.template_path
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "oci resource-manager stack create-from-template --compartment-id ${var.tenancy_ocid} --template-id ${jsondecode(data.local_file.input_template.content).data.items[0].id} --terraform-version ${var.terraform_version} --display-name ${var.resource_name} --output json > ${local.resource_path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.resource_path
}

output "display_name" {
  value = jsondecode(data.local_file.input.content).data.display-name
}

output "resource_id" {
  value = jsondecode(data.local_file.input.content).data.id
}

output "resource_name" {
  value = var.resource_name
}

output "status" {
  value = jsondecode(data.local_file.input.content).data.lifecycle-state
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}