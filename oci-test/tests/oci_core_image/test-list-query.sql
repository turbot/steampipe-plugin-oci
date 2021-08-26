select display_name, id, operating_system, operating_system_version
from oci.oci_core_image
where display_name = '{{ output.resource_name.value }}' and region = '{{ output.region.value }}';