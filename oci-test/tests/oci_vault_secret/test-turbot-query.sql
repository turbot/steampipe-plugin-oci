select title, tenant_id
from oci_vault_secret
where id = '{{ output.resource_id.value }}';
