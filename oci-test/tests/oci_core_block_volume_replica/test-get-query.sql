select id, display_name, block_volume_id
from oci.oci_core_block_volume_replica
where id = '{{ output.resource_id.value }}';