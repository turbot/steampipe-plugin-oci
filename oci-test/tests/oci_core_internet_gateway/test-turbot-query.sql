select title, tenant_id
from oci.oci_core_internet_gateway
where id = '{{ output.resource_id.value }}';