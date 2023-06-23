select name, id, lifecycle_state, project_id
from oci_devops_repository
where id = '{{ output.resource_id.value }}';
