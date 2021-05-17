select id, user_id, lifecycle_state, description
from oci.oci_identity_auth_token
where id = '{{ output.resource_id.value }}';