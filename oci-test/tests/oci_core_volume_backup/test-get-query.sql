select id, display_name, volume_id, free_form_tags, lifecycle_state, size_in_mbs
from oci.oci_core_volume_backup
where id = '{{ output.resource_id.value }}';