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
  default     = ""
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

resource "oci_core_volume" "test_volume" {
  availability_domain = var.oci_ad
  compartment_id      = var.tenancy_ocid
}

resource "oci_core_instance_configuration" "test_instance_configuration" {
  compartment_id = var.tenancy_ocid
  instance_details {
    instance_type = "compute"
    launch_details {
      availability_domain = var.oci_ad
      compartment_id      = var.tenancy_ocid
      shape               = "VM.Standard.E2.1.Micro"
      source_details {
        source_type = "image"
        image_id    = oci_core_image.test_image.id
      }
    }
  }
}

resource "oci_core_instance_pool" "test_instance_pool" {
  depends_on                = [oci_core_image.test_image]
  compartment_id            = var.tenancy_ocid
  instance_configuration_id = oci_core_instance_configuration.test_instance_configuration.id
  placement_configurations {
    availability_domain = var.oci_ad
    primary_subnet_id   = oci_core_subnet.named_test_resource.id
  }
  size = 2
}

resource "oci_autoscaling_auto_scaling_configuration" "test_auto_scaling_configuration" {
  compartment_id       = var.tenancy_ocid
  cool_down_in_seconds = "300"
  freeform_tags        = { "Department" = "Finance" }
  display_name         = var.resource_name
  is_enabled           = "true"
  policies {
    capacity {
      initial = "2"
      max     = "4"
      min     = "2"
    }
    display_name = "TFPolicy"
    policy_type  = "threshold"
    rules {
      action {
        type  = "CHANGE_COUNT_BY"
        value = "1"
      }

      display_name = var.resource_name

      metric {
        metric_type = "CPU_UTILIZATION"

        threshold {
          operator = "GT"
          value    = "90"
        }
      }
    }

    rules {
      action {
        type  = "CHANGE_COUNT_BY"
        value = "-1"
      }
      display_name = "TFScaleInRule"

      metric {
        metric_type = "CPU_UTILIZATION"

        threshold {
          operator = "LT"
          value    = "1"
        }
      }
    }
  }

  auto_scaling_resources {
    id   = oci_core_instance_pool.test_instance_pool.id
    type = "instancePool"
  }
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_autoscaling_auto_scaling_configuration.test_auto_scaling_configuration.id
}

