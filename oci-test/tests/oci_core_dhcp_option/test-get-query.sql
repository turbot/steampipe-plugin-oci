select id, display_name, lifecycle_state, vcn_id
from oci.oci_core_dhcp_option
where id = '{{ output.resource_id.value }}';