select title, tenant_id
from oci.oci_mysql_configuration
where id = '{{ output.resource_id.value }}';