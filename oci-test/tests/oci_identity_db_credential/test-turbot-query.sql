select title, tenant_id
from oci.oci_identity_db_credential
where user_id = '{{ output.user_id.value }}';