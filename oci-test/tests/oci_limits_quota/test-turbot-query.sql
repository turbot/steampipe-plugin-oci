select title, tenant_id, region
from oci.oci_limits_quota
where id = '{{ output.resource_id.value }}';