select title, tenant_id
from oci.oci_identity_db_credential
where user_id = 'ddd{{ output.user_id.value }}';