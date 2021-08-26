select display_name, id, operating_system, operating_system_version
from oci.oci_core_image_custom
where display_name = '{{ resourceName }}';