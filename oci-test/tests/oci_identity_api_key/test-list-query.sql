select key_id, user_id, fingerprint, lifecycle_state
from oci.oci_identity_api_key
where user_id = '{{ output.user_id.value }}';