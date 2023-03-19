variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
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
  default     = "us-ashburn-1"
  description = "OCI region used for the test. Does not work with default region in config, so must be defined here."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

variable "notification_topic_description" {
  type        = string
  default     = "Terraform testing resource"
  description = "The description of the topic."
}

resource "oci_ons_notification_topic" "named_test_ons_topic_resource" {
  compartment_id = var.tenancy_ocid
  name           = var.resource_name
  description    = var.notification_topic_description
  freeform_tags  = { "Name" = var.resource_name }
}

resource "oci_devops_project" "named_test_resource" {
  #Required
  name = var.resource_name
  compartment_id = var.tenancy_ocid
  notification_config {
      #Required
      topic_id = oci_ons_notification_topic.named_test_ons_topic_resource.id
  }

  #Optional
  
}

output "reporting_region" {
  value = var.region
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = oci_devops_project.named_test_resource.id
}

output "notification_topic_id" {
  value = oci_ons_notification_topic.named_test_ons_topic_resource.id
}
