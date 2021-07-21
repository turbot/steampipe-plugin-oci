select id, display_name, boot_volume_id
from oci.oci_core_boot_volume_replica
where display_name = '{{ resourceName }}';