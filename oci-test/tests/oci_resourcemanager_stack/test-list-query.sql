select display_name, id, lifecycle_state
from oci.oci_resourcemanager_stack
where display_name = '{{ output.display_name.value }}';