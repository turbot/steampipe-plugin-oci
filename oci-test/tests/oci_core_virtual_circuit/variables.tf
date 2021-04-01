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
  default     = "OCI"
  description = "OCI credentials profile used for the test. Default is to use the default profile."
}

variable "region" {
  type        = string
  default     = "ap-mumbai-1"
  description = "OCI region used for the test. Default is to use the default profile."
}

provider "oci" {
  tenancy_ocid        = var.tenancy_ocid
  config_file_profile = var.config_file_profile
}


resource "oci_core_cross_connect_group" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid

  #Optional
  customer_reference_name = var.resource_name
  display_name            = var.resource_name
}


resource "oci_core_drg" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid

  #Optional
  display_name = var.resource_name
}

resource "oci_core_virtual_circuit" "named_test_resource" {
  #Required
  compartment_id = var.tenancy_ocid
  type           = "PRIVATE"

  cross_connect_mappings {

    #Optional
    cross_connect_or_cross_connect_group_id = oci_core_cross_connect_group.named_test_resource.id
    customer_bgp_peering_ip                 = "169.254.52.246/30"
    oracle_bgp_peering_ip                   = "169.254.52.245/30"
    vlan                                    = "200"
  }

  #Optional
  bandwidth_shape_name = "10 Gbps"
  gateway_id           = oci_core_drg.named_test_resource.id

  customer_asn = "12345"
  # defined_tags = {"Operations.CostCenter"= "42"}
  display_name = var.resource_name
  # freeform_tags = {"Department"= "Finance"}
  region = var.region
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "bandwidth_shape_name" {
  value = oci_core_virtual_circuit.named_test_resource.bandwidth_shape_name
}

output "freeform_tags" {
  value = oci_core_virtual_circuit.named_test_resource.freeform_tags
}

output "resource_id" {
  value = oci_core_virtual_circuit.named_test_resource.id
}

output "display_name" {
  value = oci_core_virtual_circuit.named_test_resource.display_name
}
output "time_created" {
  value = oci_core_virtual_circuit.named_test_resource.time_created
}
output "oracle_bgp_asn" {
  value = oci_core_virtual_circuit.named_test_resource.oracle_bgp_asn
}
output "provider_service_id" {
  value = oci_core_virtual_circuit.named_test_resource.provider_service_id
}
output "provider_state" {
  value = oci_core_virtual_circuit.named_test_resource.provider_state
}
output "service_type" {
  value = oci_core_virtual_circuit.named_test_resource.service_type
}

output "state" {
  value = oci_core_virtual_circuit.named_test_resource.state
}

output "type" {
  value = oci_core_virtual_circuit.named_test_resource.type
}
output "region" {
  value = var.region
}

output "bgp_session_state" {
  value = oci_core_virtual_circuit.named_test_resource.service_type
}
