select title, tenant_id, region
from oci.oci_mysql_backup
where id = '{{ output.resource_id.value }}';