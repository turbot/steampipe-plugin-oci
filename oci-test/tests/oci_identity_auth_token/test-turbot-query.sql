select title, tenant_id
from oci.oci_identity_auth_token
where id = '{{ output.resource_id.value }}';