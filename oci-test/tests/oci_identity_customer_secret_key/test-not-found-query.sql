select id, user_id, display_name, lifecycle_state, time_created
from oci.oci_identity_customer_secret_key
where display_name = '{{ output.display_name.value }}::dummy';