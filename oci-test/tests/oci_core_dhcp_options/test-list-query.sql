select id, display_name, lifecycle_state, vcn_id
from oci.oci_core_dhcp_options
where display_name = '{{ output.display_name.value }}';