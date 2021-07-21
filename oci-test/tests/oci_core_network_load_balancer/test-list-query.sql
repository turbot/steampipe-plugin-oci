select display_name, id, lifecycle_state, subnet_id
from oci.oci_core_network_load_balancer
where display_name = '{{ resourceName }}';