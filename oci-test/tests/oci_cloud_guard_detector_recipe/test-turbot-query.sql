select title, tenant_id
from oci.oci_cloud_guard_detector_recipe
where id = '{{ output.resource_id.value }}';