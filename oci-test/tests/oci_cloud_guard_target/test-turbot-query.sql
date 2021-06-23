select title, tenant_id
from oci.oci_cloud_guard_target
where id = '{{ output.resource_id.value }}';