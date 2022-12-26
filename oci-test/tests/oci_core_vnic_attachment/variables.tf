variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_core_vcn" "test_vcn" {
  cidr_block     = "10.0.0.0/16"
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
}

resource "oci_core_subnet" "test_subnet" {
  availability_domain = var.oci_ad
  cidr_block          = "10.0.1.0/24"
  display_name        = var.resource_name
  compartment_id      = var.tenancy_ocid
  vcn_id              = oci_core_vcn.test_vcn.id
  route_table_id      = oci_core_vcn.test_vcn.default_route_table_id
  security_list_ids   = [oci_core_vcn.test_vcn.default_security_list_id]
  dhcp_options_id     = oci_core_vcn.test_vcn.default_dhcp_options_id
}

locals {
  imagePath = "${path.cwd}/image.json"
}

resource "null_resource" "test_image" {
  depends_on = [oci_core_subnet.test_subnet]
  provisioner "local-exec" {
    command = "oci compute image list --compartment-id ${var.tenancy_ocid} --all --output json > ${local.imagePath}"
  }
}

data "local_file" "image" {
  depends_on = [null_resource.test_image]
  filename   = local.imagePath
}

resource "oci_core_network_security_group" "test_network_security_group" {
  compartment_id = var.tenancy_ocid
  vcn_id         = oci_core_vcn.test_vcn.id
  display_name = var.resource_name
}

resource "oci_core_instance" "test_instance" {
  depends_on = [null_resource.test_image]
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
  display_name        = var.resource_name
  shape               = "VM.Standard2.1"

  source_details {
    source_type = "image"
    source_id   = jsondecode(data.local_file.image.content).data[0].id
  }

  create_vnic_details {
    subnet_id      = oci_core_subnet.test_subnet.id
  }
}

resource "oci_core_vnic_attachment" "named_test_resource" {
  instance_id  = oci_core_instance.test_instance.id
  display_name = "SecondaryVnicAttachment_${var.resource_name}"

  create_vnic_details {
    subnet_id                 = oci_core_subnet.test_subnet.id
    display_name              = "SecondaryVnic_${var.resource_name}"
    assign_public_ip          = true
    skip_source_dest_check    = true
    assign_private_dns_record = true
  }
}

output "resource_name" {
  value = "SecondaryVnicAttachment_${var.resource_name}"
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "instance_id" {
  depends_on = [oci_core_vnic_attachment.named_test_resource]
  value = oci_core_vnic_attachment.named_test_resource.instance_id
}

output "resource_id" {
  depends_on = [oci_core_vnic_attachment.named_test_resource]
  value = oci_core_vnic_attachment.named_test_resource.id
}

output "state" {
  depends_on = [oci_core_vnic_attachment.named_test_resource]
  value = oci_core_vnic_attachment.named_test_resource.state
}

output "vnic_id" {
  depends_on = [oci_core_vnic_attachment.named_test_resource]
  value = oci_core_vnic_attachment.named_test_resource.vnic_id
}
