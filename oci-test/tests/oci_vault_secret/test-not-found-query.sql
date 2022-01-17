select name, id, lifecycle_state
from oci_vault_secret
where id = '{{ output.resource_id.value }}nf';
