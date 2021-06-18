select title, tenant_id
from oci.oci_identity_dynamic_group
where id = '{{ output.resource_id.value }}';