select id, is_auto_tune_enabled, is_hydrated, lifecycle_state
from oci.oci_core_volume
where id = '{{ output.resource_id.value }}::dummy';