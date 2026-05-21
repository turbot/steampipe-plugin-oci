select display_name, id, lifecycle_state
from oci_waf_network_address_list
where display_name = '{{ output.resource_name.value }}';
