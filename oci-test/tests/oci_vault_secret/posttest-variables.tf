resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = "oci vault secret schedule-secret-deletion --secret-id {{ output.resource_id.value }}"
  }
}
