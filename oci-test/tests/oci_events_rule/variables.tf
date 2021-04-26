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

resource "oci_ons_notification_topic" "test_notification_topic" {
  compartment_id = var.tenancy_ocid
  name = var.resource_name
}

resource "oci_events_rule" "test_rule" {
  depends_on  = [oci_ons_notification_topic.test_notification_topic]
  actions {
      actions {
          action_type = "ONS"
          is_enabled = true
          topic_id = oci_ons_notification_topic.test_notification_topic.id
      }
  }
  compartment_id = var.tenancy_ocid
  condition =  <<EOF
  {
    "eventType": "com.oraclecloud.autoscaling.changeautoscalingconfigurationcompartment",
    "cloudEventsVersion": "0.1",
    "eventTypeVersion": "2.0",
    "source": "autoscaling",
    "eventTime": "2019-09-12T18:51:27.401Z",
    "contentType": "application/json",
    "data": {
      "compartmentId": "ocid1.compartment.oc1..unique_ID",
      "compartmentName": "example_compartment",
      "resourceId": "ocid1.autoscalingconfiguration.oc1.phx.unique_ID"
    },
    "eventID": "unique_ID",
    "extensions": {
      "compartmentId": "ocid1.compartment.oc1..unique_ID"
    }
  }
  EOF
  display_name = var.resource_name
  is_enabled = true
  description = var.resource_name
  freeform_tags = {"Department"= "Finance"}
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "freeform_tags" {
  value = oci_events_rule.test_rule.freeform_tags
}

output "resource_id" {
  value = oci_events_rule.test_rule.id
}

output "topic_id" {
  value = oci_ons_notification_topic.test_notification_topic.id
}

output "action_id" {
  value = "${element(oci_events_rule.test_rule.actions[0].actions[*].id, 0)}"
}

output "region" {
  value = var.region
}