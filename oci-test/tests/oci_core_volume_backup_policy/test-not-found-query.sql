select id, display_name freeform_tags
from oci.oci_core_volume_backup_policy
where id = '{{ output.resource_id.value }}aa';