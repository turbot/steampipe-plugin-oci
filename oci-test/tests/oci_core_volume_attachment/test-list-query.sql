select id, display_name, lifecycle_state
from oci.oci_core_volume_attachment
where display_name = '{{ resourceName }}';