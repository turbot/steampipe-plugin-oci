select boot_volume_id, display_name, id, freeform_tags, lifecycle_state, size_in_gbs
from oci.oci_core_boot_volume_backup
where id = '{{ output.resource_id.value }}';