variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "config_file_profile" {
  type        = string
  default     = "DEFAULT"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "tenancy_ocid" {
  type        = string
  default     = "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecxzjetci3yjjnjqmfkr4po3hoz4p4h2q37cyljaq"
  description = "OCID of your tenancy."
}

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

data "oci_objectstorage_namespace" "test_namespace" {
  compartment_id = var.tenancy_ocid
}

resource "oci_objectstorage_bucket" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  name           = var.resource_name
  namespace      = data.oci_objectstorage_namespace.test_namespace.namespace
}

resource "oci_objectstorage_object" "test_object" {
  bucket    = oci_objectstorage_bucket.named_test_resource.name
  content   = "test"
  namespace = data.oci_objectstorage_namespace.test_namespace.namespace
  object    = "test"
}

resource "oci_core_vcn" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
}

resource "oci_core_subnet" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  cidr_block     = "10.0.0.0/16"
  vcn_id         = oci_core_vcn.named_test_resource.id
}

resource "oci_core_image" "test_image" {
  compartment_id = var.tenancy_ocid
  display_name   = var.resource_name
  image_source_details {
    source_type    = "objectStorageTuple"
    bucket_name    = oci_objectstorage_bucket.named_test_resource.name
    namespace_name = data.oci_objectstorage_namespace.test_namespace.namespace
    object_name    = oci_objectstorage_object.test_object.object
  }
}

resource "oci_core_instance" "test_instance" {
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
  shape               = "VM.Standard.E2.1.Micro"
  source_details {
    source_id   = oci_core_image.test_image.id
    source_type = "image"
  }
  create_vnic_details {
    subnet_id = oci_core_subnet.named_test_resource.id
  }
  preserve_boot_volume = false
}

resource "oci_core_boot_volume_backup" "test_boot_volume_backup" {
  depends_on     = [oci_core_instance.test_instance]
  boot_volume_id = oci_core_instance.test_instance.boot_volume_id
  display_name   = var.resource_name
  freeform_tags  = { "Department" = "Finance" }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = oci_core_boot_volume_backup.test_boot_volume_backup.id
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "boot_volume_id" {
  value = oci_core_instance.test_instance.boot_volume_id
}
