select title, tenant_id
from oci.oci_resourcemanager_stack
where id = '{{ output.resource_id.value }}';