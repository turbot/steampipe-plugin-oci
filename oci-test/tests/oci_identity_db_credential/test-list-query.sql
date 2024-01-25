select user_id, id, lifecycle_state, description
from oci.oci_identity_db_credential
where user_id = '{{ output.user_id.value }}';