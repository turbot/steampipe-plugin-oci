select id, display_name, is_auto_tune_enabled, is_hydrated, lifecycle_state, size_in_gbs
from oci.oci_core_volume
where display_name = '{{ output.resource_id.value.split('.').pop() }}';