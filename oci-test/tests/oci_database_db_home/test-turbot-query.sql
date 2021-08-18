select title, tenant_id
from oci.oci_database_db_home
where id = '{{ output.resource_id.value }}';