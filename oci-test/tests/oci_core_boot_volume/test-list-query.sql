select id, display_name, lifecycle_state
from oci.oci_core_boot_volume
where display_name = '{{ output.resource_name.value }}';