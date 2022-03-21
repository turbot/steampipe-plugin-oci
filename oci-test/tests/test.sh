#!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 ENDCOLOR="\e[0m"


# Define your function here
run_test () {
   echo -e "${BLACK}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 > temp.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BOLDGREEN}Passed -> $1 ${ENDCOLOR}"
    echo $1 >> passed_tests.txt
   fi
  echo -e "$1" >> resource_list.txt && cat temp.txt | grep "resource_name" >> resource_list.txt && echo -e "\n\n" >> resource_list.txt
  cat temp.txt >> output.txt
  rm -rf temp.txt
 }

# output.txt - store output of each test
# failed_tests.txt - names of failed test
# passed_tests.txt names of passed test

# removes files from previous test
# rm -rf output.txt failed_tests.txt passed_tests.txt
date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt
date >> resource_list.txt


run_test oci_apigateway_api
run_test oci_core_internet_gateway
run_test oci_database_software_image
run_test oci_kms_key
run_test oci_autoscaling_auto_scaling_configuration
run_test oci_core_load_balancer
run_test oci_dns_rrset
run_test oci_kms_key_version
run_test oci_budget_alert_rule
run_test oci_core_local_peering_gateway
run_test oci_dns_tsig_key
run_test oci_kms_vault
run_test oci_budget_budget
run_test oci_core_nat_gateway
run_test oci_dns_zone
run_test oci_logging_log
run_test oci_cloud_guard_configuration
run_test oci_core_network_load_balancer
run_test oci_events_rule
run_test oci_logging_log_group
run_test oci_cloud_guard_detector_recipe
run_test oci_core_network_security_group
run_test oci_file_storage_file_system
run_test oci_mysql_backup
run_test oci_cloud_guard_managed_list
run_test oci_core_public_ip
run_test oci_file_storage_mount_target
run_test oci_mysql_channel
run_test oci_cloud_guard_responder_recipe
run_test oci_core_public_ip_pool
run_test oci_file_storage_snapshot
run_test oci_mysql_configuration
run_test oci_cloud_guard_target
run_test oci_core_route_table
run_test oci_functions_application
run_test oci_mysql_configuration_custom
run_test oci_core_block_volume_replica
run_test oci_core_security_list
run_test oci_identity_api_key
run_test oci_mysql_db_system
run_test oci_core_boot_volume
run_test oci_core_service_gateway
run_test oci_identity_auth_token
run_test oci_nosql_table
run_test oci_core_boot_volume_attachment
run_test oci_core_subnet
run_test oci_identity_availability_domain
run_test oci_objectstorage_bucket
run_test oci_core_boot_volume_backup
run_test oci_core_vcn
run_test oci_identity_customer_secret_key
run_test oci_ons_notification_topic
run_test oci_core_boot_volume_replica
run_test oci_core_vnic_attachment
run_test oci_identity_dynamic_group
run_test oci_ons_subscription
run_test oci_core_dhcp_options
run_test oci_core_volume
run_test oci_identity_network_source
run_test oci_resource_search
run_test oci_core_drg
run_test oci_core_volume_attachment
run_test oci_identity_policy
run_test oci_vault_secret
run_test oci_core_image
run_test oci_core_volume_backup
run_test oci_identity_tag_default
run_test oci_core_image_custom
run_test oci_core_volume_backup_policy
run_test oci_identity_tag_namespace
run_test oci_core_instance
run_test oci_database_autonomous_database
run_test oci_identity_tenancy

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt