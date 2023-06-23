select name, id, lifecycle_state
from oci_devops_repository
where name = '{{ resourceName }}';
