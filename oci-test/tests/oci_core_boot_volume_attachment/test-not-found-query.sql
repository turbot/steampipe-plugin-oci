select id
from oci.oci_core_boot_volume_attachment
where id = '{{ output.resource_id.value }}::dummy';