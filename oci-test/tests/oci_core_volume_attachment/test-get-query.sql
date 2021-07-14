select id, is_read_only, lifecycle_state
from oci.oci_core_volume_attachment
where id = '{{ output.resource_id.value }}';