select title, tenant_id, region
from oci_core_service_gateway
where id = '{{ output.resource_id.value }}aa';