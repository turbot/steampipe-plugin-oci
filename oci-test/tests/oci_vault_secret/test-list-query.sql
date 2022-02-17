select name, id, lifecycle_state
from oci_vault_secret
where name = '{{ output.resource_name.value }}';
