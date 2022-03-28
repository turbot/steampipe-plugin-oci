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

variable "template_id" {
  type        = string
  default     = "ocid1.ormtemplate.oc1.ap-mumbai-1.aaaaaaaa2ereypft5pb3vgmplr6arm767dpqvkfocugw4vdcbrgmmg2dbcsq"
  description = "OCID of template."
}

variable "terraform-version" {
  type        = string
  default     = "0.12.x"
  description = "The version of Terraform to use with the stack."
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
    command = "oci resource-manager stack create-from-template --compartment-id ${var.tenancy_ocid} --template-id ${var.template_id} --terraform-version ${var.terraform-version} --output json > ${local.path} --display-name ${var.resource_name}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
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