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
  default     = ""
  description = "OCID of your tenancy."
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

resource "oci_database_db_home" "named_test_resource" {
  #Required
  database {
    admin_password = var.db_home_database_admin_password

    #Optional
    backup_id                  = oci_database_backup.test_backup.id
    backup_tde_password        = var.db_home_database_backup_tde_password
    character_set              = var.db_home_database_character_set
    database_id                = oci_database_database.test_database.id
    database_software_image_id = oci_database_database_software_image.test_database_software_image.id
    db_backup_config {

      #Optional
      auto_backup_enabled = var.db_home_database_db_backup_config_auto_backup_enabled
      auto_backup_window  = var.db_home_database_db_backup_config_auto_backup_window
      backup_destination_details {

        #Optional
        id   = var.db_home_database_db_backup_config_backup_destination_details_id
        type = var.db_home_database_db_backup_config_backup_destination_details_type
      }
      recovery_window_in_days = var.db_home_database_db_backup_config_recovery_window_in_days
    }
    db_name                               = var.db_home_database_db_name
    db_workload                           = var.db_home_database_db_workload
    defined_tags                          = var.db_home_database_defined_tags
    freeform_tags                         = var.db_home_database_freeform_tags
    ncharacter_set                        = var.db_home_database_ncharacter_set
    pdb_name                              = var.db_home_database_pdb_name
    tde_wallet_password                   = var.db_home_database_tde_wallet_password
    time_stamp_for_point_in_time_recovery = var.db_home_database_time_stamp_for_point_in_time_recovery
  }

  #Optional
  database_software_image_id = oci_database_database_software_image.test_database_software_image.id
  db_system_id               = oci_database_db_system.test_db_system.id
  db_version {
  }
  defined_tags           = var.db_home_defined_tags
  display_name           = var.db_home_display_name
  freeform_tags          = { "Department" = "Finance" }
  is_desupported_version = var.db_home_is_desupported_version
  kms_key_id             = oci_kms_key.test_key.id
  kms_key_version_id     = oci_kms_key_version.test_key_version.id
  source                 = var.db_home_source
  vm_cluster_id          = oci_database_vm_cluster.test_vm_cluster.id
}

output "resource_name" {
  value = var.resource_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_database_db_home.named_test_resource.id
}
