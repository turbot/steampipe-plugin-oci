select title, tenant_id
from oci.oci_database_database
where id = '{{ output.resource_id.value }}';