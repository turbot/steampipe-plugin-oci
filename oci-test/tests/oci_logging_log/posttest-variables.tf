resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = <<EOT
    oci logging log delete --log-group-id {{ output.log_group_id.value }} --log-id {{ output.resource_id.value }} --force
    EOT
  }
}
