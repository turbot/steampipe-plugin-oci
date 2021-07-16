select id, display_name, is_auto_tune_enabled, is_hydrated, lifecycle_state, size_in_gbs
from oci.oci_core_boot_volume
where display_name = '{{ resourceName }}';