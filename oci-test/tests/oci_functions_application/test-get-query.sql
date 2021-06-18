select id, display_name, lifecycle_state, subnet_ids, freeform_tags
from oci.oci_functions_application
where id = '{{ output.resource_id.value }}';