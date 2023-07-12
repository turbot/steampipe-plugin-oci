select name, id
from oci_devops_repository
where id = 'demo-test-{{ output.resource_id.value }}';
