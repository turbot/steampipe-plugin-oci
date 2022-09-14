select title, tenant_id
from oci.oci_containerengine_cluster
where id = '{{ output.resource_id.value }}';