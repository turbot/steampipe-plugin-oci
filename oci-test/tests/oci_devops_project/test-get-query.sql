select id, name, freeform_tags, lifecycle_state
from oci.oci_devops_project
where id = '{{ output.resource_id.value }}';