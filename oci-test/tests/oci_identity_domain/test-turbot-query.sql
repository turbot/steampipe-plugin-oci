select title, tenant_id
from oci.oci_identity_domain
where id = '{{ output.resource_id.value }}';