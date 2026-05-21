select display_name, id, lifecycle_state
from oci_waf_network_address_list
where id = '{{ output.resource_id.value }}';
