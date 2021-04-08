select display_name, id, volume_id, freeform_tags, lifecycle_state, size_in_gbs
from oci.oci_core_volume_backup
where display_name = '{{ resourceName }}';