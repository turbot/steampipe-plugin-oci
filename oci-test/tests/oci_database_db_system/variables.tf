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

resource "oci_database_db_system" "test_db_system" {
  #Required
  availability_domain = var.db_system_availability_domain
  compartment_id      = var.compartment_id
  source              = "NONE"
  db_home {
    database {
      #Required
      admin_password = var.db_system_db_home_database_admin_password

      #Optional
      character_set              = var.db_system_db_home_database_character_set
      database_software_image_id = oci_database_database_software_image.test_database_software_image.id
      db_backup_config {

        #Optional
        auto_backup_enabled = false
        backup_destination_details {
          #Optional
          id   = var.db_system_db_home_database_db_backup_config_backup_destination_details_id
          type = var.db_system_db_home_database_db_backup_config_backup_destination_details_type
        }
        recovery_window_in_days = var.db_system_db_home_database_db_backup_config_recovery_window_in_days
      }
      db_domain                             = var.db_system_db_home_database_db_domain
      db_name                               = var.db_system_db_home_database_db_name
      db_workload                           = var.db_system_db_home_database_db_workload
      defined_tags                          = var.db_system_db_home_database_defined_tags
      freeform_tags                         = var.db_system_db_home_database_freeform_tags
      ncharacter_set                        = var.db_system_db_home_database_ncharacter_set
      pdb_name                              = var.db_system_db_home_database_pdb_name
      tde_wallet_password                   = var.db_system_db_home_database_tde_wallet_password
      time_stamp_for_point_in_time_recovery = var.db_system_db_home_database_time_stamp_for_point_in_time_recovery
    }

    #Optional
    database_software_image_id = oci_database_database_software_image.test_database_software_image.id
    db_version                 = var.db_system_db_home_db_version
    defined_tags               = var.db_system_db_home_defined_tags
    display_name               = var.db_system_db_home_display_name
    freeform_tags              = var.db_system_db_home_freeform_tags
  }
  hostname        = var.db_system_hostname
  shape           = var.db_system_shape
  ssh_public_keys = var.db_system_ssh_public_keys
  subnet_id       = oci_core_subnet.test_subnet.id

  #Optional
  backup_network_nsg_ids  = var.db_system_backup_network_nsg_ids
  backup_subnet_id        = oci_core_subnet.test_subnet.id
  cluster_name            = var.db_system_cluster_name
  cpu_core_count          = var.db_system_cpu_core_count
  data_storage_percentage = var.db_system_data_storage_percentage
  data_storage_size_in_gb = var.db_system_data_storage_size_in_gb
  database_edition        = var.db_system_database_edition
  db_system_options {

    #Optional
    storage_management = var.db_system_db_system_options_storage_management
  }
  defined_tags       = var.db_system_defined_tags
  disk_redundancy    = var.db_system_disk_redundancy
  display_name       = var.db_system_display_name
  domain             = var.db_system_domain
  fault_domains      = var.db_system_fault_domains
  freeform_tags      = { "Department" = "Finance" }
  kms_key_id         = oci_kms_key.test_key.id
  kms_key_version_id = oci_kms_key_version.test_key_version.id
  license_model      = var.db_system_license_model
  maintenance_window_details {

    #Optional
    days_of_week {

      #Optional
      name = var.db_system_maintenance_window_details_days_of_week_name
    }
    hours_of_day       = var.db_system_maintenance_window_details_hours_of_day
    lead_time_in_weeks = var.db_system_maintenance_window_details_lead_time_in_weeks
    months {

      #Optional
      name = var.db_system_maintenance_window_details_months_name
    }
    preference     = var.db_system_maintenance_window_details_preference
    weeks_of_month = var.db_system_maintenance_window_details_weeks_of_month
  }
  node_count          = var.db_system_node_count
  nsg_ids             = var.db_system_nsg_ids
  private_ip          = var.db_system_private_ip
  source_db_system_id = oci_database_db_system.test_db_system.id
  sparse_diskgroup    = var.db_system_sparse_diskgroup
  time_zone           = var.db_system_time_zone
}

output "resource_name" {
  value = var.resource_name
}

output "db_name" {
  value = oci_database_db_system.named_test_resource.db_name
}

output "tenancy_ocid" {
  value = var.tenancy_ocid
}

output "resource_id" {
  value = oci_database_db_system.named_test_resource.id
}
