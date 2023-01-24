select id, display_name, lifecycle_state, size_in_gbs
from oci.oci_core_volume_group
where display_name = '{{ output.display_name.value }}';