select id, display_name, is_enabled, lifecycle_state
from oci.oci_core_internet_gateway
where id = '{{ output.resource_id.value }}';