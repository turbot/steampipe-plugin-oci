select title, tenant_id
from oci.oci_cloud_guard_responder_recipe
where id = '{{ output.resource_id.value }}';