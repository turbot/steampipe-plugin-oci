select title, tenant_id
from oci.oci_core_local_peering_gateway
where id = '{{ output.resource_id.value }}';