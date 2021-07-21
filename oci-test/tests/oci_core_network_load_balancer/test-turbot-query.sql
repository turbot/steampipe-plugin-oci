select title, tenant_id, region
from oci.oci_core_network_load_balancer
where id = '{{ output.resource_id.value }}';