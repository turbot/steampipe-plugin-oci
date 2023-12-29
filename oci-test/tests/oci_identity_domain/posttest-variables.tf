resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = "oci iam domain deactivate --domain-id {{ output.resource_id.value }}"
  }
  provisioner "local-exec" {
    command = "sleep 30"
  }
  provisioner "local-exec" {
    command = "oci iam domain delete --domain-id {{ output.resource_id.value }} --force"
  }
}