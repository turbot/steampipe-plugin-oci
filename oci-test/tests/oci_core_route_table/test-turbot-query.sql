select title, tenant_id
from oci.oci_core_route_table
where id = '{{ output.resource_id.value }}';