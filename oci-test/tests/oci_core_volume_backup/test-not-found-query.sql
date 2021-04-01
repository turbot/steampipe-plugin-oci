select id, display_name, lifecycle_state, free_form_tags
from oci.oci_core_volume_backup
where id = '{{ output.resource_id.value }}aa';