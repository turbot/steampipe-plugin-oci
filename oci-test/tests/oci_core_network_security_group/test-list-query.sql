select display_name, id, lifecycle_state
from oci.oci_core_network_security_group
where display_name = '{{ resourceName }}';