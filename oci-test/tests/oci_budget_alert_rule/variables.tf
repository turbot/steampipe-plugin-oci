variable "resource_name" {
  type        = string
  default     = "steampipetest20200125"
  description = "Name of the resource used throughout the test."
}

variable "tenancy_ocid" {
  type        = string
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

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}

resource "oci_budget_budget" "named_test_resource" {
  display_name          = var.resource_name
  amount                = 100
  compartment_id        = var.tenancy_ocid
  reset_period          = "MONTHLY"
  target_compartment_id = var.tenancy_ocid
}

resource "oci_budget_alert_rule" "named_test_resource" {
  budget_id      = oci_budget_budget.named_test_resource.id
  display_name   = var.resource_name
  threshold      = "100"
  threshold_type = "PERCENTAGE"
  type           = "ACTUAL"
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_budget_alert_rule.named_test_resource.id
}

output "budget_id" {
  value = oci_budget_budget.named_test_resource.id
}
