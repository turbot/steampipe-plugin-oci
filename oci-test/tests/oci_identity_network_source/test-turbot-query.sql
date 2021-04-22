select title, tenant_id
from oci.oci_identity_network_source
where id = '{{ output.resource_id.value }}';