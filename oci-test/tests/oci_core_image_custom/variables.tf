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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

variable "image" {
  type        = string
  default     = "Oracle-Autonomous-Linux-7.9-2021.05-0"
  description = "Oracle supported platform image."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_core_vcn" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
}

resource "oci_core_subnet" "named_test_resource" {
  depends_on     = [oci_core_vcn.named_test_resource]
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
  vcn_id         = oci_core_vcn.named_test_resource.id
  freeform_tags  = { "Name" = var.resource_name }
}

locals {
  imagePath = "${path.cwd}/image.json"
  path      = "${path.cwd}/output.json"
}

resource "null_resource" "test_image" {
  depends_on = [oci_core_subnet.named_test_resource]
  provisioner "local-exec" {
    command = "oci compute image list --compartment-id ${var.tenancy_ocid} --all --display-name ${var.image} --output json > ${local.imagePath}"
  }
}

data "local_file" "image" {
  depends_on = [null_resource.test_image]
  filename   = local.imagePath
}

resource "null_resource" "named_test_resource" {
  depends_on = [null_resource.test_image]
  provisioner "local-exec" {
    command = "oci compute instance launch --availability-domain ${var.oci_ad} --compartment-id ${var.tenancy_ocid} --shape VM.Standard2.1 --subnet-id ${oci_core_subnet.named_test_resource.id} --image-id ${jsondecode(data.local_file.image.content).data[0].id} --output json > ${local.path}"
  }
  provisioner "local-exec" {
    command = "sleep 150"
  }
}

data "local_file" "instance" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

resource "oci_core_image" "named_test_resource" {
  depends_on = [null_resource.named_test_resource]
    compartment_id = var.tenancy_ocid
    instance_id = jsondecode(data.local_file.instance.content).data.id
    display_name = var.resource_name
}

resource "null_resource" "destroy_test_resource" {
  depends_on = [oci_core_image.named_test_resource]
  provisioner "local-exec" {
    command = "oci compute instance terminate --instance-id ${jsondecode(data.local_file.instance.content).data.id} --force"
  }
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_core_image.named_test_resource.id
}

output "lifecycle_state" {
  value = oci_core_image.named_test_resource.state
}

output "operating_system" {
  value = oci_core_image.named_test_resource.operating_system
}

output "operating_system_version" {
  value = oci_core_image.named_test_resource.operating_system_version
}
