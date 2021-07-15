select title, tenant_id, region
from oci.oci_mysql_channel
where id = '{{ output.resource_id.value }}';