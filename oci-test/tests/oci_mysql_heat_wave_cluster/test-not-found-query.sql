select db_system_id, lifecycle_state, time_created
from oci.oci_mysql_heat_wave_cluster
where db_system_id = '{{ output.resource_id.value }}aa';