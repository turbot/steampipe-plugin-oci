select display_name, id, lifecycle_state, description, db_system_id, backup_type, creation_type, mysql_version
from oci.oci_mysql_db_backup
where display_name = '{{ resourceName }}';