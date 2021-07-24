select id, display_name, boot_volume_id
from oci.oci_core_boot_volume_replica
where boot_volume_id = '{{ output.boot_volume_id.value }}' and lifecycle_state = 'AVAILABLE';