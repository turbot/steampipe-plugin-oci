select title, tenant_id, region
from oci.oci_events_rule
where id = '{{ output.resource_id.value }}';