select id, user_id, display_name, lifecycle_state, title
from oci.oci_identity_customer_secret_key
where id  = '{{ output.resource_id.value }}';