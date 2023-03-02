select title, tenant_id, region
from oci.oci_queue_queue
where id = '{{ output.resource_id.value }}';
