select title, tenant_id
from oci.oci_identity_customer_secret_key
where display_name = '{{ output.display_name.value }}';