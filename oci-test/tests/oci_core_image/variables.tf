variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
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

variable "oci_ad" {
  type        = string
  default     = "TvRS:AP-MUMBAI-1-AD-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "local_file" "foo" {
    content     = "foo!"
    filename = "${path.module}/foo.bar"
}

resource "oci_objectstorage_bucket" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  name = var.resource_name
  namespace = "bmqeqvslavsz"
}

resource "oci_objectstorage_object" "test_object" {
  bucket = oci_objectstorage_bucket.named_test_resource.name
  content = local_file.foo.filename
  namespace = "bmqeqvslavsz"
  object = local_file.foo.filename
}

resource "oci_core_image" "named_test_resource" {
    compartment_id = var.tenancy_ocid
    display_name = var.resource_name
    freeform_tags = {"Name"= var.resource_name}
    image_source_details {
      source_type = "objectStorageTuple"
      bucket_name = oci_objectstorage_bucket.named_test_resource.name
      namespace_name = "bmqeqvslavsz"
      object_name = oci_objectstorage_object.test_object.object
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


