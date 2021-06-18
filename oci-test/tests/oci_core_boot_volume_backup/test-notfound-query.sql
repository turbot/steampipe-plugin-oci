select id, display_name, lifecycle_state, freeform_tags
from oci.oci_core_boot_volume_backup
where id = '{{ output.resource_id.value }}aa';