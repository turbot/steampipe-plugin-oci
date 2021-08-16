resource "null_resource" "destroy_test_resource" {
  provisioner "local-exec" {
    command = "oci bv boot-volume update --boot-volume-id {{ output.boot_volume_id.value }} --boot-volume-replicas '[]' --force"
  }
  provisioner "local-exec" {
    command = "oci compute instance terminate --instance-id {{ output.instance_id.value }} --force"
  }
}
