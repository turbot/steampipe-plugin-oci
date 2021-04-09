select id, display_name, lifecycle_state
from oci.oci_core_internet_gateway
where id = '{{ output.resource_id.value }}::dummy';