select id, name
from oci.oci_devops_project
where id = '{{ output.resource_id.value }}aa';