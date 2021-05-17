select display_name, id, freeform_tags, region
from oci.oci_core_volume_backup_policy
where display_name = '{{ resourceName }}';