resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = "oci mysql configuration delete --configuration-id {{ output.resource_id.value }} --force"
  }
}
