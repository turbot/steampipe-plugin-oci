select title, region, tenant_id
from oci.oci_logging_log
where name = '{{ output.display_name.value }}';