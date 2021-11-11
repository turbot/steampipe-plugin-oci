variable "user_name" {
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
  region              = var.region
}

resource "oci_identity_user" "named_test_resource" {
  compartment_id = var.tenancy_ocid
  description    = var.user_name
  name           = var.user_name
}

resource "oci_identity_api_key" "named_test_resource" {
  depends_on = [oci_identity_user.named_test_resource]
  key_value = <<EOF
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtBLQAGmKJ7tpfzYJyqLG
ZDwHL51+d6T8Z00BnP9CFfzxZZZ48PcYSUHuTyCM8mR5JqYLyH6C8tZ/DKqwxUnc
ONgBytG3MM42bgxfHIhsZRj5rCz1oqWlSLuXvgww1kuqWnt6r+NtnXog439YsGTH
RotrTLTdEgOxH0EFP5uHUc9w/Uix7rWU7GB2ra060oeTB/hKpts5U70eI2EI6ec9
1sJdUIj7xNfBJeQQrz4CFUrkyzL06211CFvhmxH2hA9gBKOqC3rGL8XraHZBhGWn
mXlrQB7nNKsJrrv5fHwaPDrAY4iNP2W0q3LRpyNigJ6cgRuGJhHa82iHPmxgIx8m
fwIDAQAB
-----END PUBLIC KEY-----
EOF
  user_id = oci_identity_user.named_test_resource.id
  provisioner "local-exec" {
    command = "sleep 120"
  }
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "user_id" {
  depends_on = [oci_identity_user.named_test_resource]
  value = oci_identity_api_key.named_test_resource.user_id
}

output "key_id" {
  depends_on = [oci_identity_api_key.named_test_resource]
  value = oci_identity_api_key.named_test_resource.id
}

output "fingerprint" {
  depends_on = [oci_identity_api_key.named_test_resource]
  value = oci_identity_api_key.named_test_resource.fingerprint
}

output "state" {
  depends_on = [oci_identity_api_key.named_test_resource]
  value = oci_identity_api_key.named_test_resource.state
}

