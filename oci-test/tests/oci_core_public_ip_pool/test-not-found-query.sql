select title, tenant_id
from oci.oci_core_public_ip_pool
where id = '{{ output.resource_id.value }}::dummy';