resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = "oci compute instance terminate --instance-id {{ output.resource_id.value }} --force"
  }
}
