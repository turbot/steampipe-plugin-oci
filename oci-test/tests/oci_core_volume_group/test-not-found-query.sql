select id, is_hydrated, lifecycle_state
from oci.oci_core_volume_group
where id = '{{ output.resource_id.value }}::dummy';