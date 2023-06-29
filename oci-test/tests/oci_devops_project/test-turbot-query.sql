select title, tenant_id
from oci.oci_devops_project
where id = '{{ output.resource_id.value }}';