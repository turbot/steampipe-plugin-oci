select id, lifecycle_state
from oci.oci_core_block_volume_replica
where id = '{{ output.resource_id.value }}::dummy';