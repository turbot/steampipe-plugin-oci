select title, tenant_id
from oci.oci_identity_tag_namespace
where id = '{{ output.resource_id.value }}';