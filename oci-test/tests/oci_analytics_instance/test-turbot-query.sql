select title, tenant_id
from oci.oci_analytics_instance
where id = '{{ output.resource_id.value }}';