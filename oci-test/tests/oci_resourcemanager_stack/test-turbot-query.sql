select title, tenant_id, region
from oci.oci_resourcemanager_stack
where id = '{{ output.resource_id.value }}';