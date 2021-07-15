select id, display_name, lifecycle_state
from oci.oci_core_block_volume_replica
where display_name = '{{ resourceName }}';