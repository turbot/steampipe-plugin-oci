select title, tenant_id
from oci.oci_identity_policy
where id = '{{ output.resource_id.value }}';