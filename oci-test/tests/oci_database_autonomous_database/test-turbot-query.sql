select title, tenant_id
from oci.oci_database_autonomous_database
where id = '{{ output.resource_id.value }}';