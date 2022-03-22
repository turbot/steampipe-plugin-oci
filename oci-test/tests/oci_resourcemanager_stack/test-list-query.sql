select display_name, id, lifecycle_state, time_created
from oci.oci_resourcemanager_stack
where display_name = '{{ resourceName }}';