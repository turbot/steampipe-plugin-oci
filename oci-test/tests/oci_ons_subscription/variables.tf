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

variable "notification_topic_description" {
  type        = string
  default     = "Terraform testing resource"
  description = "The description of the topic."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
  region              = var.region
}

resource "oci_ons_notification_topic" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  name           = var.resource_name
  description    = var.notification_topic_description
}

resource "oci_ons_subscription" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  endpoint       = "test@gmail.com"
  protocol       = "EMAIL"
  topic_id       = oci_ons_notification_topic.named_test_resource.topic_id
  freeform_tags  = { "Name" = var.resource_name }
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "region" {
  value = var.region
}

output "resource_id" {
  value = oci_ons_subscription.named_test_resource.id
}

output "created_time" {
  value = oci_ons_subscription.named_test_resource.created_time
}

output "endpoint" {
  value = oci_ons_subscription.named_test_resource.endpoint
}

output "etag" {
  value = oci_ons_subscription.named_test_resource.etag
}

output "protocol" {
  value = oci_ons_subscription.named_test_resource.protocol
}

output "lifecycle_state" {
  value = oci_ons_subscription.named_test_resource.state
}
