select title, tenant_id
from oci.oci_identity_api_key
where user_id = '{{ output.user_id.value }}';