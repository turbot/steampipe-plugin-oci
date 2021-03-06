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
  default     = "ap-hyderabad-1"
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
  imagePath    = "${path.cwd}/image.json"
  path         = "${path.cwd}/output.json"
  instancePath = "${path.cwd}/instance.json"
  volumePath   = "${path.cwd}/volume.json"
}

resource "null_resource" "test_image" {
  depends_on = [oci_core_subnet.named_test_resource]
  provisioner "local-exec" {
    command = "oci compute image list --compartment-id ${var.tenancy_ocid} --all --output json > ${local.imagePath}"
  }
}

data "local_file" "image" {
  depends_on = [null_resource.test_image]
  filename   = local.imagePath
}

resource "null_resource" "test_instance" {
  depends_on = [null_resource.test_image]
  provisioner "local-exec" {
    command = "oci compute instance launch --availability-domain ${var.oci_ad} --compartment-id ${var.tenancy_ocid} --shape VM.Standard2.1 --subnet-id ${oci_core_subnet.named_test_resource.id} --image-id ${jsondecode(data.local_file.image.content).data[0].id} --output json > ${local.instancePath}"
  }
  provisioner "local-exec" {
    command = "sleep 120"
  }
}

data "local_file" "instance" {
  depends_on = [null_resource.test_instance]
  filename   = local.instancePath
}

resource "null_resource" "volume_resource" {
  provisioner "local-exec" {
    command = "oci compute boot-volume-attachment get --boot-volume-attachment-id ${jsondecode(data.local_file.instance.content).data.id} > ${local.volumePath}"
  }
}

data "local_file" "volume" {
  depends_on = [null_resource.volume_resource]
  filename   = local.volumePath
}

resource "null_resource" "named_test_resource" {
  depends_on = [null_resource.volume_resource]
  provisioner "local-exec" {
    command = "oci bv boot-volume update --boot-volume-id ${jsondecode(data.local_file.volume.content).data.boot-volume-id} --boot-volume-replicas '[{\"displayName\":\"test\",\"availabilityDomain\":\"TvRS:AP-HYDERABAD-1-AD-1\"}]' --force > ${local.path}"
  }
  provisioner "local-exec" {
    command = "sleep 120"
  }
}

data "local_file" "replice_volume" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "instance_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.instance.content).data.id
}

output "resource_name" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.replice_volume.content).data.boot-volume-replicas[0].display-name
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.replice_volume.content).data.boot-volume-replicas[0].boot-volume-replica-id
}

output "boot_volume_id" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.volume.content).data.boot-volume-id
}
