select title, tenant_id
from oci.oci_nosql_table
where id = '{{ output.resource_id.value }}';