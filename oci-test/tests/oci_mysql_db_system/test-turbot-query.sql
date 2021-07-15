select title, tenant_id
from oci.oci_mysql_db_system
where id = '{{ output.resource_id.value }}';