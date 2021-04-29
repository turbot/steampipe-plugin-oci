select title, region, tenant_id
from oci.oci_logging_log_group
where display_name = '{{ output.display_name.value }}';