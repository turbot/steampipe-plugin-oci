select id, display_name, volume_id
from oci.oci_core_volume_attachment
where display_name = '{{ resourceName }}';