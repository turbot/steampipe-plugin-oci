select title, tenant_id
from oci_devops_repository
where id = '{{ output.resource_id.value }}';