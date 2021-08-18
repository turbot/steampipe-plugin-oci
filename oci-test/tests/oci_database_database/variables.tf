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

resource "oci_database_database" "named_test_resource" {
  #Required
  database {
    admin_password = var.database_database_admin_password
    db_name        = var.database_database_db_name

    #Optional
    backup_id                  = oci_database_backup.test_backup.id
    backup_tde_password        = var.database_database_backup_tde_password
    character_set              = var.database_database_character_set
    database_software_image_id = oci_database_database_software_image.test_database_software_image.id
    db_backup_config {

      #Optional
      auto_backup_enabled = var.database_database_db_backup_config_auto_backup_enabled
      auto_backup_window  = var.database_database_db_backup_config_auto_backup_window
      backup_destination_details {

        #Optional
        id   = var.database_database_db_backup_config_backup_destination_details_id
        type = var.database_database_db_backup_config_backup_destination_details_type
      }
      recovery_window_in_days = var.database_database_db_backup_config_recovery_window_in_days
    }
    db_unique_name      = var.database_database_db_unique_name
    db_workload         = var.database_database_db_workload
    defined_tags        = var.database_database_defined_tags
    freeform_tags       = var.database_database_freeform_tags
    ncharacter_set      = var.database_database_ncharacter_set
    pdb_name            = var.database_database_pdb_name
    tde_wallet_password = var.database_database_tde_wallet_password
  }
  db_home_id = oci_database_db_home.test_db_home.id
  source     = var.database_source

  #Optional
  db_version         = var.database_db_version
  kms_key_id         = oci_kms_key.test_key.id
  kms_key_version_id = oci_kms_key_version.test_key_version.id
}

output "resource_name" {
  value = var.resource_name
}

output "db_name" {
  value = oci_database_database.named_test_resource.db_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_database_database.named_test_resource.id
}
